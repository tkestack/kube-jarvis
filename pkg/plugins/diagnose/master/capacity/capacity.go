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
package capacity

import (
	"context"
	"fmt"

	"tkestack.io/kube-jarvis/pkg/util"

	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "master-capacity"
)

// Capacity define a health master node resource status
type Capacity struct {
	// Memory is total memory of master node
	Memory resource.Quantity
	// Cpu is total core number of master node
	Cpu resource.Quantity
	// MaxNodeTotal indicate the max node number of this master scale
	MaxNodeTotal int
}

// Diagnostic check whether the resources are sufficient for a specific size cluster
type Diagnostic struct {
	*diagnose.MetaData
	result     chan *diagnose.Result
	Capacities []Capacity
	param      *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a master-node diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:   make(chan *diagnose.Result, 100),
		MetaData: meta,
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	if len(d.Capacities) == 0 {
		d.Capacities = DefCapacities
	}

	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 100)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		d.diagnoseCapacity(ctx)
	}()
	return d.result, nil
}

func (d *Diagnostic) diagnoseCapacity(ctx context.Context) {
	nodeTotal := 0
	masters := make([]v12.Node, 0)
	for _, n := range d.param.Resources.Nodes.Items {
		for k := range n.Labels {
			if k == "node-role.kubernetes.io/master" {
				masters = append(masters, n)
				continue
			}
		}
		nodeTotal++
	}

	scale, err := d.targetCapacity(nodeTotal)
	if err != nil {
		d.result <- &diagnose.Result{
			Level:   diagnose.HealthyLevelFailed,
			ObjName: "*",
			Title:   "Failed",
			Desc:    translate.Message(err.Error()),
		}
		return
	}

	for _, m := range masters {
		cpu := m.Status.Capacity.Cpu()
		mem := m.Status.Capacity.Memory()
		if cpu.Cmp(scale.Cpu) < 0 {
			d.sendCapacityWarnResult(m.Name, "Cpu", nodeTotal, util.CpuQuantityStr(cpu), util.CpuQuantityStr(&scale.Cpu))
		} else {
			d.sendCapacityGoodResult(m.Name, "Cpu", nodeTotal, util.CpuQuantityStr(cpu), util.CpuQuantityStr(&scale.Cpu))
		}

		if mem.Cmp(scale.Memory) < 0 {
			d.sendCapacityWarnResult(m.Name, "Memory", nodeTotal, util.MemQuantityStr(mem), util.MemQuantityStr(&scale.Memory))
		} else {
			d.sendCapacityGoodResult(m.Name, "Memory", nodeTotal, util.MemQuantityStr(mem), util.MemQuantityStr(&scale.Memory))
		}
	}
}

func (d *Diagnostic) sendCapacityWarnResult(name string, resource string, nTotal int, curVal, targetVal string) {
	objInfo := map[string]interface{}{
		"NodeName":    name,
		"Resource":    resource,
		"TargetValue": targetVal,
		"CurValue":    curVal,
		"NodeTotal":   nTotal,
	}

	d.result <- &diagnose.Result{
		ObjName: name,
		Level:   diagnose.HealthyLevelWarn,
		ObjInfo: objInfo,
		Title: d.Translator.Message("title", map[string]interface{}{
			"Resource": resource,
		}),
		Desc:     d.Translator.Message("desc", objInfo),
		Proposal: d.Translator.Message("proposal", objInfo),
	}
}

func (d *Diagnostic) sendCapacityGoodResult(name string, resource string, nTotal int, curVal, targetVal string) {
	objInfo := map[string]interface{}{
		"NodeName":    name,
		"Resource":    resource,
		"TargetValue": targetVal,
		"CurValue":    curVal,
		"NodeTotal":   nTotal,
	}

	d.result <- &diagnose.Result{
		ObjName: name,
		Level:   diagnose.HealthyLevelGood,
		ObjInfo: objInfo,
		Title: d.Translator.Message("title", map[string]interface{}{
			"Resource": resource,
		}),
		Desc: d.Translator.Message("good-desc", objInfo),
	}
}

func (d *Diagnostic) targetCapacity(nTotal int) (*Capacity, error) {
	for _, scale := range d.Capacities {
		if scale.MaxNodeTotal > nTotal {
			return &scale, nil
		}
	}
	return nil, fmt.Errorf("no target capacity found")
}
