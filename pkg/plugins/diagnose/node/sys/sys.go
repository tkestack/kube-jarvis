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
package sys

import (
	"context"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "node-sys"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a example diagnostic
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
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)

	go func() {
		defer diagnose.CommonDeafer(d.result)
		for node := range d.param.Resources.Machines {
			// net.ipv4.tcp_tw_reuse
			d.diagnoseKernelParam("net.ipv4.tcp_tw_reuse", "1", node)

			// net.ipv4.ip_forward
			d.diagnoseKernelParam("net.ipv4.ip_forward", "1", node)

			// net.bridge.bridge-nf-call-iptables
			d.diagnoseKernelParam("net.bridge.bridge-nf-call-iptables", "1", node)
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) diagnoseKernelParam(key string, targetVal string, node string) {
	m := d.param.Resources.Machines[node]
	curVal := m.SysCtl[key]
	level := diagnose.HealthyLevelGood

	obj := map[string]interface{}{
		"Node":      node,
		"Name":      key,
		"CurVal":    curVal,
		"TargetVal": targetVal,
	}

	if curVal != targetVal {
		level = diagnose.HealthyLevelWarn
		d.result <- &diagnose.Result{
			Level:    level,
			Title:    d.Translator.Message("kernel-para-title", nil),
			ObjName:  node,
			ObjInfo:  obj,
			Desc:     d.Translator.Message("kernel-para-desc", obj),
			Proposal: d.Translator.Message("kernel-para-proposal", obj),
		}
	} else {
		d.result <- &diagnose.Result{
			Level:   level,
			Title:   d.Translator.Message("kernel-para-title", nil),
			ObjName: node,
			Desc:    d.Translator.Message("kernel-para-good-desc", obj),
		}
	}
}
