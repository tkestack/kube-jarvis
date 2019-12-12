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

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	go func() {
		defer close(d.result)
		for node, m := range d.param.Resources.Machines {
			curVal := m.SysCtl["net.ipv4.tcp_tw_reuse"]
			targetVal := "1"
			level := diagnose.HealthyLevelGood
			if curVal != targetVal {
				level = diagnose.HealthyLevelWarn
			}

			d.result <- &diagnose.Result{
				Level:   level,
				Title:   d.Translator.Message("kernel-para-title", nil),
				ObjName: node,
				Desc: d.Translator.Message("kernel-para-desc", map[string]interface{}{
					"Node":      node,
					"Name":      "net.ipv4.tcp_tw_reuse",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),

				Proposal: d.Translator.Message("kernel-para-proposal", map[string]interface{}{
					"Node":      node,
					"Name":      "net.ipv4.tcp_tw_reuse",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),
			}
		}
	}()
	return d.result, nil
}
