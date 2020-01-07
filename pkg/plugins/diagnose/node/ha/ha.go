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
package ha

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType   = "node-ha"
	FailureDomainKey = "failure-domain.beta.kubernetes.io/zone"
)

// Diagnostic is a ha diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
}

// NewDiagnostic return a ha diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:   make(chan *diagnose.Result, 1000),
		MetaData: meta,
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		nodes := param.Resources.Nodes
		if nodes != nil {
			d.checkNodeNum(nodes)
			d.checkNodeZone(nodes)
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) checkNodeNum(nodes *v1.NodeList) {
	objName := "node-num"
	num := len(nodes.Items)
	info := map[string]interface{}{
		"Name":         objName,
		"CurTotalNode": num,
	}
	if len(nodes.Items) <= 1 {
		d.sendResult(diagnose.HealthyLevelSerious, objName, "bad", info)
	} else {
		d.sendResult(diagnose.HealthyLevelGood, objName, "good", info)
	}
}

func (d *Diagnostic) checkNodeZone(nodes *v1.NodeList) {
	if len(nodes.Items) <= 1 {
		return
	}
	objName := "node-zone"
	zoneCpu := make(map[string]int64)
	zoneMemory := make(map[string]int64)
	info := map[string]interface{}{
		"Name":            objName,
		"CurTotalZoneNum": 0,
		"ZoneName":        "",
		"ResourceName":    "",
	}
	var totalCpu, totalMemory, maxZoneCpu, maxZoneMemory int64
	for i := 0; i < len(nodes.Items); i++ {
		if nodes.Items[i].Labels != nil {
			id, ok := nodes.Items[i].Labels[FailureDomainKey]
			if !ok {
				d.sendFailedResult(objName, fmt.Errorf("lack of failure domain key:%s", FailureDomainKey))
				return
			}
			cpu, _ := nodes.Items[i].Status.Allocatable.Cpu().AsInt64()
			memory, _ := nodes.Items[i].Status.Allocatable.Memory().AsInt64()
			zoneCpu[id] += cpu
			zoneMemory[id] += memory
			totalCpu += cpu
			totalMemory += memory
		}
	}
	info["CurTotalZoneNum"] = len(zoneCpu)
	for zoneName, cpu := range zoneCpu {
		if cpu > maxZoneCpu {
			maxZoneCpu = cpu
			info["ZoneName"] = zoneName
		}
	}
	for zoneName, mem := range zoneMemory {
		if mem > maxZoneMemory {
			maxZoneMemory = mem
			info["ZoneName"] = zoneName
		}
	}
	if len(zoneCpu) == 1 {
		info["ResourceName"] = "zone"
		d.sendResult(diagnose.HealthyLevelWarn, objName, "bad", info)
	} else if (totalCpu-maxZoneCpu >= maxZoneCpu) && (totalMemory-maxZoneMemory >= maxZoneMemory) {
		d.sendResult(diagnose.HealthyLevelGood, objName, "good", info)
	} else if totalCpu-maxZoneCpu < maxZoneCpu {
		info["ResourceName"] = "cpu"
		d.sendResult(diagnose.HealthyLevelWarn, objName, "bad", info)
	} else if totalMemory-maxZoneMemory < maxZoneMemory {
		info["ResourceName"] = "memory"
		d.sendResult(diagnose.HealthyLevelWarn, objName, "bad", info)
	}

}

func (d *Diagnostic) sendResult(level diagnose.HealthyLevel, objName, descType string, extra map[string]interface{}) {

	d.result <- &diagnose.Result{
		Level:    level,
		ObjName:  objName,
		ObjInfo:  extra,
		Title:    d.Translator.Message(objName+"-title", nil),
		Desc:     d.Translator.Message(objName+"-"+descType+"-desc", extra),
		Proposal: d.Translator.Message(objName+"-proposal", extra),
	}
}

func (d *Diagnostic) sendFailedResult(objName string, err error) {

	d.result <- &diagnose.Result{
		Level:   diagnose.HealthyLevelFailed,
		ObjName: "*",
		ObjInfo:  map[string]interface{}{
			"ResourceName":"none",
			"CurTotalZoneNum":0,
		},
		Title:   "Failed",
		Desc:    translate.Message(err.Error()),
	}
}
