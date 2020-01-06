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
package iptables

import (
	"context"

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "node-iptables"

	GoodIPTablesCount = 100
	WarnIPTablesCount = 6000
	RiskIPTablesCount = 10000
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
	d.result = make(chan *diagnose.Result, 1000)
	d.param = &param
	go func() {
		defer diagnose.CommonDeafer(d.result)
		for node, m := range d.param.Resources.Machines {
			var cntLevel diagnose.HealthyLevel
			totalCount := m.IPTables.Filter.Count + m.IPTables.NAT.Count
			if totalCount < GoodIPTablesCount {
				cntLevel = diagnose.HealthyLevelGood
			} else if totalCount < WarnIPTablesCount {
				cntLevel = diagnose.HealthyLevelWarn
			} else if totalCount < RiskIPTablesCount {
				cntLevel = diagnose.HealthyLevelRisk
			} else {
				cntLevel = diagnose.HealthyLevelSerious
			}
			if cntLevel != diagnose.HealthyLevelGood {
				obj := map[string]interface{}{
					"Node":           node,
					"Name":           "iptables-count",
					"SuggestedCount": GoodIPTablesCount,
					"CurCount":       totalCount,
				}

				d.result <- &diagnose.Result{
					Level:    cntLevel,
					Title:    d.Translator.Message("iptables-count-title", nil),
					ObjName:  node,
					ObjInfo:  obj,
					Desc:     d.Translator.Message("iptables-count-desc", obj),
					Proposal: d.Translator.Message("iptables-count-proposal", obj),
				}
			}

			obj := map[string]interface{}{
				"Node":            node,
				"Name":            "iptables-forward-policy",
				"CurPolicy":       m.IPTables.Filter.ForwardPolicy,
				"SuggestedPolicy": cluster.AcceptPolicy,
			}

			forwardPolicyLevel := diagnose.HealthyLevelGood
			if m.IPTables.Filter.ForwardPolicy != cluster.AcceptPolicy {
				forwardPolicyLevel = diagnose.HealthyLevelWarn
				d.result <- &diagnose.Result{
					Level:    forwardPolicyLevel,
					Title:    d.Translator.Message("iptables-forward-policy-title", nil),
					ObjName:  node,
					ObjInfo:  obj,
					Desc:     d.Translator.Message("iptables-forward-policy-desc", obj),
					Proposal: d.Translator.Message("iptables-forward-policy-proposal", obj),
				}
			} else {
				d.result <- &diagnose.Result{
					Level:   forwardPolicyLevel,
					Title:   d.Translator.Message("iptables-forward-policy-title", nil),
					ObjName: node,
					ObjInfo: obj,
					Desc:    d.Translator.Message("iptables-forward-policy-good-desc", obj),
				}
			}
		}
	}()
	return d.result, nil
}
