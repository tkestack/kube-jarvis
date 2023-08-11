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
package components

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"time"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "master-components"

	FailureDomainKey = "failure-domain.beta.kubernetes.io/zone"
)

// Diagnostic check that the core components are working properly (include k8s node components)
//also check if they have been restarted within 24 hours
type Diagnostic struct {
	*diagnose.MetaData
	result      chan *diagnose.Result
	Components  []string
	RestartTime string
}

// NewDiagnostic return a master-components diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:   make(chan *diagnose.Result, 1000),
		MetaData: meta,
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	if d.RestartTime != "" {
		_, err := time.ParseDuration(d.RestartTime)
		if err != nil {
			return errors.Wrapf(err, "wrong config : restarttime=%s : %v", d.RestartTime, err)
		}
	} else {
		d.RestartTime = "24h"
	}

	if len(d.Components) == 0 {
		d.Components = []string{
			cluster.ComponentApiserver,
			cluster.ComponentScheduler,
			cluster.ComponentControllerManager,
			cluster.ComponentETCD,
			cluster.ComponentKubeProxy,
			cluster.ComponentCoreDNS,
			cluster.ComponentKubeDNS,
			cluster.ComponentKubelet,
			cluster.ComponentDockerd,
			cluster.ComponentContainerd,
		}
	}

	return nil
}

func isMasterCoreComp(comp string) bool {
	return comp == cluster.ComponentApiserver ||
		comp == cluster.ComponentScheduler ||
		comp == cluster.ComponentControllerManager ||
		comp == cluster.ComponentETCD
}

func getCompResultLevel(comp string) diagnose.HealthyLevel {
	if isMasterCoreComp(comp) {
		return diagnose.HealthyLevelSerious
	}
	return diagnose.HealthyLevelRisk
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context,
	param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	_, err := time.ParseDuration(d.RestartTime)
	if err != nil {
		return nil, errors.Wrapf(err, "wrong config : restarttime=%s : %v", d.RestartTime, err)
	}

	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		etcdNodes := make([]v1.Node, 0)
		var etcdReplica int
		for _, comp := range d.Components {
			compInfos, exist := param.Resources.CoreComponents[comp]
			if !exist {
				d.sendCompNotExist(comp)
				return
			}
			for _, inf := range compInfos {
				if inf.Name == cluster.ComponentETCD {
					etcdReplica += 1
					for _, n := range param.Resources.Nodes.Items {
						node := n
						if node.Name == inf.Node {
							etcdNodes = append(etcdNodes, node)
							break
						}
					}
				}
				if inf.Error != nil {
					d.sendNormalResult(comp, &inf, diagnose.HealthyLevelFailed, "err", map[string]interface{}{
						"Err": inf.Error.Error(),
					})
					continue
				}

				if !inf.IsRunning {
					d.sendNormalResult(comp, &inf, getCompResultLevel(comp), "not-run", nil)
					continue
				}

				if had, exTra := d.hadRestart(inf.Pod); had {
					d.sendNormalResult(comp, &inf, diagnose.HealthyLevelRisk, "restart", exTra)
				} else {
					d.sendNormalResult(comp, &inf, diagnose.HealthyLevelGood, "good", nil)
				}
			}
		}
		zone := make(map[string]int)
		for _, node := range etcdNodes {
			if node.Labels != nil {
				id, ok := node.Labels[FailureDomainKey]
				if !ok {
					d.Logger.Errorf("lack of failure domain key: %s, node: %s", FailureDomainKey, node.Name)
					return
				}
				zone[id] += 1
			}
		}
		if len(zone) < 3 {
			d.sendAbnormalResult(diagnose.HealthyLevelRisk, "ha", etcdReplica, len(zone))
		}

	}()
	return d.result, nil
}

func (d *Diagnostic) hadRestart(p *v1.Pod) (bool, map[string]interface{}) {
	if p == nil {
		return false, nil
	}

	dt, _ := time.ParseDuration(d.RestartTime)
	for _, s := range p.Status.ContainerStatuses {
		if s.RestartCount > 0 &&
			s.State.Running != nil &&
			s.State.Running.StartedAt.Add(dt).After(time.Now()) {
			return true, map[string]interface{}{
				"Count":    s.RestartCount,
				"LastTime": s.State.Running.StartedAt.String(),
			}
		}
	}

	return false, nil
}

func (d *Diagnostic) sendCompNotExist(comp string) {
	d.result <- &diagnose.Result{
		Level:   diagnose.HealthyLevelFailed,
		ObjName: comp,
		Title:   "Failed",
		Desc:    "can not found target component info",
	}
}

func (d *Diagnostic) sendNormalResult(comp string, inf *cluster.Component,
	level diagnose.HealthyLevel, preFix string, extra map[string]interface{}) {
	obj := map[string]interface{}{
		"Name":      inf.Name,
		"Node":      inf.Node,
		"Component": comp,
	}

	for k, v := range extra {
		obj[k] = v
	}

	d.result <- &diagnose.Result{
		Level:    level,
		ObjName:  inf.Name,
		ObjInfo:  obj,
		Title:    d.Translator.Message(preFix+"-title", obj),
		Desc:     d.Translator.Message(preFix+"-desc", obj),
		Proposal: d.Translator.Message(preFix+"-proposal", obj),
	}
}

func (d *Diagnostic) sendAbnormalResult(level diagnose.HealthyLevel, preFix string, replica, curTotalZoneNum int) {
	obj := map[string]interface{}{
		"Name":            "node-zone",
		"Component":       "etcd",
		"Replica":         fmt.Sprint(replica),
		"CurTotalZoneNum": fmt.Sprint(curTotalZoneNum),
	}

	d.result <- &diagnose.Result{
		Level:    level,
		ObjName:  "etcd-node-zone",
		ObjInfo:  obj,
		Title:    d.Translator.Message(preFix+"-title", obj),
		Desc:     d.Translator.Message(preFix+"-desc", obj),
		Proposal: d.Translator.Message(preFix+"-proposal", obj),
	}
}
