package custom

import (
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"sync"
	"tkestack.io/kube-jarvis/pkg/logger"
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
		Components: map[string]*compexplorer.Auto{
			cluster.ComponentApiserver:         compexplorer.NewAuto(cluster.ComponentApiserver, true),
			cluster.ComponentScheduler:         compexplorer.NewAuto(cluster.ComponentScheduler, true),
			cluster.ComponentControllerManager: compexplorer.NewAuto(cluster.ComponentControllerManager, true),
			cluster.ComponentETCD:              compexplorer.NewAuto(cluster.ComponentETCD, true),
			cluster.ComponentKubeProxy:         compexplorer.NewAuto(cluster.ComponentKubeProxy, false),
			cluster.ComponentCoreDNS:           compexplorer.NewAuto(cluster.ComponentCoreDNS, false),
			cluster.ComponentKubeDNS:           compexplorer.NewAuto(cluster.ComponentKubeDNS, false),
			cluster.ComponentKubelet:           compexplorer.NewAuto(cluster.ComponentKubelet, false),
			cluster.ComponentDockerd:           compexplorer.NewAuto(cluster.ComponentKubelet, false),
			cluster.ComponentContainerd:        compexplorer.NewAuto(cluster.ComponentContainerd, false),
		},
	}

	return c
}

// Init do initialization for Cluster
func (c *Cluster) Init() error {
	var err error
	c.nodeExecutor, err = c.Node.Executor(c.logger, c.cli, c.restConfig)
	if err != nil && err != nodeexec.NoneExecutor {
		return errors.Wrap(err, "create node executor failed")
	}

	for t, cmp := range c.Components {
		if err := cmp.Init(c.logger, c.cli, c.nodeExecutor); err != nil {
			return errors.Wrapf(err, "init component executor for %s failed", t)
		}
		c.compExps[t] = cmp
	}

	return nil
}

// SyncResources fetch all resource from cluster
func (c *Cluster) SyncResources() error {
	c.resources = cluster.NewResources()
	if err := c.initK8sResources(); err != nil {
		return err
	}

	if err := c.initComponents(); err != nil {
		return err
	}

	return c.initMachines()
}

func (c *Cluster) initK8sResources() (err error) {
	c.logger.Infof("Fetching all k8s resources..")
	client := c.cli.CoreV1()
	admissionControllerClient := c.cli.AdmissionregistrationV1beta1()
	opts := v1.ListOptions{}

	var g errgroup.Group
	g.Go(func() (err error) {
		c.resources.Nodes, err = client.Nodes().List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.PersistentVolumes, err = client.PersistentVolumes().List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.ComponentStatuses, err = client.ComponentStatuses().List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.Pods, err = client.Pods(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.PodTemplates, err = client.PodTemplates(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.PersistentVolumeClaims, err = client.PersistentVolumeClaims(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.ConfigMaps, err = client.ConfigMaps(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.Secrets, err = client.Secrets(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.Services, err = client.Services(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.ServiceAccounts, err = client.ServiceAccounts(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.ResourceQuotas, err = client.ResourceQuotas(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.LimitRanges, err = client.LimitRanges(v1.NamespaceAll).List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.SystemNamespace, err = client.Namespaces().Get(v1.NamespaceSystem, v1.GetOptions{})
		return
	})
	g.Go(func() (err error) {
		c.resources.MutatingWebhookConfigurations, err = admissionControllerClient.MutatingWebhookConfigurations().List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.ValidatingWebhookConfigurations, err = admissionControllerClient.ValidatingWebhookConfigurations().List(opts)
		return
	})
	g.Go(func() (err error) {
		c.resources.Namespaces, err = client.Namespaces().List(opts)
		return
	})

	return g.Wait()
}

func (c *Cluster) initComponents() error {
	for name, cmp := range c.compExps {
		result, err := cmp.Component()
		if err != nil {
			return errors.Wrapf(err, "fetch component %s failed ", name)
		}

		c.resLock.Lock()
		c.resources.CoreComponents[name] = result
		c.resLock.Unlock()
		c.logger.Infof("Fetching (%d) %s ..", len(result), name)
	}

	return nil
}

func (c *Cluster) initMachines() error {
	nodes, err := c.cli.CoreV1().Nodes().List(v1.ListOptions{})
	if err != nil {
		return errors.Wrapf(err, "get nodes list failed")
	}
	c.logger.Infof("Fetching (%d) machines information ..", len(nodes.Items))

	conCtl := make(chan struct{}, 100)
	var g errgroup.Group
	for _, n := range nodes.Items {
		node := n
		g.Go(func() error {
			conCtl <- struct{}{}
			defer func() { <-conCtl }()

			m, err := c.getOneNodeInfo(node.Name)
			if err != nil {
				return errors.Wrapf(err, "get machine info failed")
			}

			c.resLock.Lock()
			c.resources.Machines[node.Name] = m
			c.resLock.Unlock()
			return nil
		})
	}

	return g.Wait()
}

func (c *Cluster) getOneNodeInfo(nodeName string) (cluster.Machine, error) {
	out, _, err := c.nodeExecutor.DoCmd(nodeName, []string{"sysctl", "-a"})
	if err != nil {
		return cluster.Machine{}, errors.Wrapf(err, "do commond 'sysctl -a' failed")
	}

	return cluster.Machine{
		SysCtl: GetSysCtlMap(out),
	}, nil
}

func (c *Cluster) Resources() *cluster.Resources {
	return c.resources
}

// CloudType return the cloud type of Cluster
func (c *Cluster) CloudType() string {
	return Type
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
