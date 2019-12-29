/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package custom

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/compexplorer"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

const (
	// Type is the cluster type
	Type = "custom"
)

// Cluster is a custom cluster
type Cluster struct {
	Node       *nodeexec.Config
	Components map[string]*compexplorer.Auto
	KubeConfig string

	cli          kubernetes.Interface
	restConfig   *rest.Config
	logger       logger.Logger
	nodeExecutor nodeexec.Executor
	resources    *cluster.Resources
	resLock      sync.Mutex

	compExps map[string]compexplorer.Explorer
	progress *plugins.Progress
}

// NewCluster return an new custom Cluster
func NewCluster(log logger.Logger, cli kubernetes.Interface, config *rest.Config) cluster.Cluster {
	c := &Cluster{
		logger:     log,
		cli:        cli,
		restConfig: config,
		Node:       nodeexec.NewConfig(),
		resources:  cluster.NewResources(),
		compExps:   map[string]compexplorer.Explorer{},
		Components: map[string]*compexplorer.Auto{},
	}

	return c
}

// Complete check and complete config items
func (c *Cluster) Complete() error {
	if _, exist := c.Components[cluster.ComponentApiserver]; !exist {
		c.Components[cluster.ComponentApiserver] = compexplorer.NewAuto(cluster.ComponentApiserver, true)
	}

	if _, exist := c.Components[cluster.ComponentScheduler]; !exist {
		c.Components[cluster.ComponentScheduler] = compexplorer.NewAuto(cluster.ComponentScheduler, true)
	}

	if _, exist := c.Components[cluster.ComponentControllerManager]; !exist {
		c.Components[cluster.ComponentControllerManager] = compexplorer.NewAuto(cluster.ComponentControllerManager, true)
	}
	if _, exist := c.Components[cluster.ComponentETCD]; !exist {
		c.Components[cluster.ComponentETCD] = compexplorer.NewAuto(cluster.ComponentETCD, true)
	}
	if _, exist := c.Components[cluster.ComponentKubeProxy]; !exist {
		c.Components[cluster.ComponentKubeProxy] = compexplorer.NewAuto(cluster.ComponentKubeProxy, false)
	}

	if _, exist := c.Components[cluster.ComponentCoreDNS]; !exist {
		c.Components[cluster.ComponentCoreDNS] = compexplorer.NewAuto(cluster.ComponentCoreDNS, false)
	}

	if _, exist := c.Components[cluster.ComponentKubeDNS]; !exist {
		c.Components[cluster.ComponentKubeDNS] = compexplorer.NewAuto(cluster.ComponentKubeDNS, false)
	}

	if _, exist := c.Components[cluster.ComponentKubelet]; !exist {
		c.Components[cluster.ComponentKubelet] = compexplorer.NewAuto(cluster.ComponentKubelet, false)
	}

	if _, exist := c.Components[cluster.ComponentDockerd]; !exist {
		c.Components[cluster.ComponentDockerd] = compexplorer.NewAuto(cluster.ComponentDockerd, false)
	}

	if _, exist := c.Components[cluster.ComponentContainerd]; !exist {
		c.Components[cluster.ComponentContainerd] = compexplorer.NewAuto(cluster.ComponentContainerd, false)
	}

	for _, cmp := range c.Components {
		cmp.Complete()
	}

	c.Node.Complete()

	return nil
}

// Init do initialization for Cluster
func (c *Cluster) Init(ctx context.Context, progress *plugins.Progress) error {
	c.progress = progress
	c.resources = cluster.NewResources()
	c.progress.CreateStep("init_env", "Preparing environment", 2)
	c.progress.CreateStep("init_k8s_resources", "Fetching k8s resources..", 20)
	c.progress.CreateStep("init_components", "Fetching all components..", len(c.compExps))
	nodes, err := c.cli.CoreV1().Nodes().List(v1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "get nodes from k8s failed")
	}
	c.progress.CreateStep("init_machines", "Fetching all machines..", len(nodes.Items))

	// now start init steps
	c.logger.Infof("Start preparing environment...........")
	c.progress.SetCurStep("init_env")
	if err := c.initExecutors("init_env"); err != nil {
		return err
	}

	c.logger.Infof("Start fetching all k8s resources...........")
	c.progress.SetCurStep("init_k8s_resources")
	if err := c.initK8sResources("init_k8s_resources"); err != nil {
		return err
	}

	c.logger.Infof("Start fetching all components...........")
	c.progress.SetCurStep("init_components")
	if err := c.initComponents("init_components"); err != nil {
		return err
	}

	c.logger.Infof("Start fetching all machines...........")
	c.progress.SetCurStep("init_machines")
	if err := c.initMachines("init_machines"); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) initExecutors(stepName string) error {
	var err error
	if c.nodeExecutor == nil {
		c.nodeExecutor, err = c.Node.Executor(c.logger, c.cli, c.restConfig)
		if err != nil && err != nodeexec.NoneExecutor {
			return errors.Wrap(err, "create node executor failed")
		}
	}
	c.progress.AddStepPercent(stepName, 1)

	for t, cmp := range c.Components {
		if err := cmp.Init(c.logger, c.cli, c.nodeExecutor); err != nil {
			return errors.Wrapf(err, "init component executor for %s failed", t)
		}
		if c.compExps[t] == nil {
			c.compExps[t] = cmp
		}
	}
	c.progress.AddStepPercent(stepName, 1)

	return nil
}

func (c *Cluster) initK8sResources(stepName string) error {
	client := c.cli.CoreV1()
	admissionControllerClient := c.cli.AdmissionregistrationV1beta1()
	opts := v1.ListOptions{}
	var g errgroup.Group
	g.Go(func() (err error) {
		c.resources.Nodes, err = client.Nodes().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list Nodes failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Nodes", len(c.resources.Nodes.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.PersistentVolumes, err = client.PersistentVolumes().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list PersistentVolumes failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) PersistentVolumes", len(c.resources.PersistentVolumes.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.ComponentStatuses, err = client.ComponentStatuses().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list ComponentStatuses failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) ComponentStatuses", len(c.resources.ComponentStatuses.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Pods, err = client.Pods(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list Pods failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Pods", len(c.resources.Pods.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.PodTemplates, err = client.PodTemplates(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list PodTemplates failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) PodTemplates", len(c.resources.PodTemplates.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.PersistentVolumeClaims, err = client.PersistentVolumeClaims(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list PersistentVolumeClaims failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) PersistentVolumeClaims", len(c.resources.PersistentVolumeClaims.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.ConfigMaps, err = client.ConfigMaps(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list ConfigMaps failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) ConfigMaps", len(c.resources.ConfigMaps.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Secrets, err = client.Secrets(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list Secrets failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Secrets", len(c.resources.Secrets.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Services, err = client.Services(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list Services failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Services", len(c.resources.Services.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.ServiceAccounts, err = client.ServiceAccounts(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list ServiceAccounts failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) ServiceAccounts", len(c.resources.ServiceAccounts.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.ResourceQuotas, err = client.ResourceQuotas(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list ResourceQuotas failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) ResourceQuotas", len(c.resources.ResourceQuotas.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.LimitRanges, err = client.LimitRanges(v1.NamespaceAll).List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list LimitRanges failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) LimitRanges", len(c.resources.LimitRanges.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.MutatingWebhookConfigurations, err = admissionControllerClient.MutatingWebhookConfigurations().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list MutatingWebhookConfigurations failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) MutatingWebhookConfigurations", len(c.resources.MutatingWebhookConfigurations.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.ValidatingWebhookConfigurations, err = admissionControllerClient.ValidatingWebhookConfigurations().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list ValidatingWebhookConfigurations failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) ValidatingWebhookConfigurations", len(c.resources.ValidatingWebhookConfigurations.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Namespaces, err = client.Namespaces().List(opts)
		if err != nil {
			err = errors.Wrapf(err, "list Namespaces failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Namespaces", len(c.resources.Namespaces.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Deployments, err = c.cli.AppsV1().Deployments("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list Deployments failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Deployments", len(c.resources.Deployments.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.DaemonSets, err = c.cli.AppsV1().DaemonSets("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list DaemonSets failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) DaemonSets", len(c.resources.DaemonSets.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.StatefulSets, err = c.cli.AppsV1().StatefulSets("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list StatefulSets failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) StatefulSets", len(c.resources.StatefulSets.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.Jobs, err = c.cli.BatchV1().Jobs("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list Jobs failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) Jobs", len(c.resources.Jobs.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.CronJobs, err = c.cli.BatchV1beta1().CronJobs("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list CronJobs failed")
		} else {
			c.progress.AddStepPercent(stepName, 1)
			c.logger.Infof("Fetching (%d) CronJobs", len(c.resources.CronJobs.Items))
		}
		return
	})

	g.Go(func() (err error) {
		c.resources.HPAs, err = c.cli.AutoscalingV1().HorizontalPodAutoscalers("").List(v1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "list HPAs failed")
		}
		c.logger.Infof("Fetching (%d) HPAs", len(c.resources.HPAs.Items))
		return
	})

	return g.Wait()
}

func (c *Cluster) initComponents(stepName string) error {
	g := errgroup.Group{}
	for tempName, tempCmp := range c.compExps {
		name := tempName
		cmp := tempCmp
		g.Go(func() error {
			result, err := cmp.Component()
			if err != nil {
				return errors.Wrapf(err, "fetch component %s failed ", name)
			}
			c.logger.Infof("Fetching (%d) %s", len(result), name)
			c.resLock.Lock()
			c.resources.CoreComponents[name] = result
			c.resLock.Unlock()
			c.progress.AddStepPercent(stepName, 1)
			return nil
		})
	}

	return g.Wait()
}

func (c *Cluster) initMachines(stepName string) error {
	nodes, err := c.cli.CoreV1().Nodes().List(v1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "get nodes from k8s failed")
	}

	var g errgroup.Group
	conCtl := make(chan struct{}, 200)
	for _, n := range nodes.Items {
		node := n
		g.Go(func() error {
			conCtl <- struct{}{}
			defer func() { <-conCtl }()

			m := c.getOneNodeInfo(node.Name)
			c.resLock.Lock()
			c.resources.Machines[node.Name] = m
			c.resLock.Unlock()

			c.progress.AddStepPercent(stepName, 1)
			return nil
		})
	}

	return g.Wait()
}

func (c *Cluster) getOneNodeInfo(nodeName string) cluster.Machine {
	out, errStr, err := c.nodeExecutor.DoCmd(nodeName, []string{"sh", "-c", "sysctl -a | grep -v error"})
	if err != nil {
		c.logger.Errorf("Failed to get node %s sysctl set: %s, %v", nodeName, errStr, err)
		return cluster.Machine{
			Error: errors.Wrapf(err, "do commond 'sysctl -a' failed"),
		}
	}

	out1, errStr, err := c.nodeExecutor.DoCmd(nodeName, []string{"sh", "-c", "iptables-save"})
	if err != nil {
		c.logger.Errorf("Failed to get node %s iptables info: %s, %v", nodeName, errStr, err)
		return cluster.Machine{
			Error: errors.Wrapf(err, "do commond 'iptables-save' failed"),
		}
	}

	sysctlSet := GetSysCtlMap(out)
	//c.logger.Debugf("Get node %s sysctl result: %v", nodeName, sysctlSet)
	iptablesInfo := GetIPTablesInfo(out1)
	c.logger.Debugf("Get node %s iptables result: %v", nodeName, iptablesInfo)
	return cluster.Machine{
		SysCtl:   sysctlSet,
		IPTables: iptablesInfo,
	}
}

func (c *Cluster) Resources() *cluster.Resources {
	return c.resources
}

// CloudType return the cloud type of Cluster
func (c *Cluster) CloudType() string {
	return Type
}

// Finish will be called once diagnostic done
func (c *Cluster) Finish() error {
	if err := c.nodeExecutor.Finish(); err != nil {
		return errors.Wrapf(err, "finish node executor failed")
	}
	c.nodeExecutor = nil

	for t, cmp := range c.compExps {
		if err := cmp.Finish(); err != nil {
			return errors.Wrapf(err, "finis component explore %s failed", t)
		}
	}
	c.compExps = map[string]compexplorer.Explorer{}
	return nil
}

// GetSysCtlMap cover a out from "sysctl -a" to a map
func GetSysCtlMap(out string) map[string]string {
	result := map[string]string{}
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		spIndex := strings.IndexAny(line, "=")
		if spIndex == -1 {
			continue
		}

		k := line[0:spIndex]
		v := line[spIndex+1:]
		result[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return result
}

// GetIPTablesInfo cover a out from "iptables-save"
func GetIPTablesInfo(out string) (result cluster.IPTablesInfo) {
	lines := strings.Split(out, "\n")

	var idx int
	result.NAT, idx = getNATTableInfo(lines)
	result.Filter = getFilterTableInfo(lines[idx+1:])
	return
}

func getNATTableInfo(lines []string) (nat cluster.NATTable, end int) {
	end = -1
	var found bool
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "*nat" {
			found = true
		}
		if !found {
			continue
		}
		nat.Count++
		if line == "COMMIT" {
			end = idx
			return
		}
		if strings.HasPrefix(line, ":PREROUTING ") {
			if strings.Contains(line, "ACCEPT") {
				nat.PreRoutingPolicy = cluster.AcceptPolicy
			} else {
				nat.PreRoutingPolicy = cluster.DropPolicy
			}
		}
		if strings.HasPrefix(line, ":INPUT ") {
			if strings.Contains(line, "ACCEPT") {
				nat.InputPolicy = cluster.AcceptPolicy
			} else {
				nat.InputPolicy = cluster.DropPolicy
			}
		}
		if strings.HasPrefix(line, ":OUTPUT ") {
			if strings.Contains(line, "ACCEPT") {
				nat.OutputPolicy = cluster.AcceptPolicy
			} else {
				nat.OutputPolicy = cluster.DropPolicy
			}
		}
		if strings.HasPrefix(line, ":POSTROUTING ") {
			if strings.Contains(line, "ACCEPT") {
				nat.PostRoutingPolicy = cluster.AcceptPolicy
			} else {
				nat.PostRoutingPolicy = cluster.DropPolicy
			}
		}
	}
	return
}

func getFilterTableInfo(lines []string) (filter cluster.FilterTable) {
	var found bool
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if line == "*filter" {
			found = true
		}
		if !found {
			continue
		}
		filter.Count++
		if line == "COMMIT" {
			return
		}
		if strings.HasPrefix(line, ":INPUT ") {
			if strings.Contains(line, "ACCEPT") {
				filter.InputPolicy = cluster.AcceptPolicy
			} else {
				filter.InputPolicy = cluster.DropPolicy
			}
		}
		if strings.HasPrefix(line, ":FORWARD ") {
			if strings.Contains(line, "ACCEPT") {
				filter.ForwardPolicy = cluster.AcceptPolicy
			} else {
				filter.ForwardPolicy = cluster.DropPolicy
			}
		}
		if strings.HasPrefix(line, ":OUTPUT ") {
			if strings.Contains(line, "ACCEPT") {
				filter.OutputPolicy = cluster.AcceptPolicy
			} else {
				filter.OutputPolicy = cluster.DropPolicy
			}
		}
	}
	return
}
