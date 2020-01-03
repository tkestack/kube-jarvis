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
	"math"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType       = "node-ha"
	FailureDomainKey     = "failure-domain.beta.kubernetes.io/zone"
	DefaultZoneNodeRatio = 0.6
)

// Diagnostic is a ha diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result        chan *diagnose.Result
	ZoneNodeRatio float64 `yaml: "zoneNodeRatio"`
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
	if math.Abs(d.ZoneNodeRatio-0) < 1e-3 {
		d.ZoneNodeRatio = DefaultZoneNodeRatio
	}
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
	zone := make(map[string]int)
	info := map[string]interface{}{
		"Name":            objName,
		"CurTotalZoneNum": 0,
		"CurNodeRatio":    0.0,
		"ErrMsg":          "",
	}
	var id string
	var ok, fail bool
	for i := 0; i < len(nodes.Items); i++ {
		if nodes.Items[i].Labels != nil {
			id, ok = nodes.Items[i].Labels[FailureDomainKey]
			if !ok {
				info["ErrMsg"] = fmt.Sprintf("lack of failure domain key:%s", FailureDomainKey)
				d.sendResult(diagnose.HealthyLevelFailed, objName, "bad", info)
				fail = true
				break
			}
			zone[id]++
		}
	}
	if fail {
		return
	}
	info["CurTotalZoneNum"] = len(zone)
	if len(zone) == 1 {
		d.sendResult(diagnose.HealthyLevelWarn, objName, "bad", info)
	} else if len(zone) == 2 {
		var num int
		var ratio float64
		for _, value := range zone {
			if num == 0 {
				num = value
			} else if value <= num {
				ratio = float64(value * 1.0 / num)
			} else if value > num {
				ratio = float64(num * 1.0 / value)
			}
		}
		info["CurNodeRatio"] = ratio
		if ratio >= d.ZoneNodeRatio {
			d.sendResult(diagnose.HealthyLevelGood, objName, "good", info)
		} else {
			d.sendResult(diagnose.HealthyLevelWarn, objName, "bad", info)
		}
	} else {
		d.sendResult(diagnose.HealthyLevelGood, objName, "good", info)
	}

}

func (d *Diagnostic) sendResult(level diagnose.HealthyLevel, objName, descType string, extra map[string]interface{}) {

	d.result <- &diagnose.Result{
		Level:    level,
		ObjName:  objName,
		Title:    d.Translator.Message(objName+"-title", nil),
		Desc:     d.Translator.Message(objName+"-"+descType+"-desc", extra),
		Proposal: d.Translator.Message(objName+"-proposal", extra),
	}
}
