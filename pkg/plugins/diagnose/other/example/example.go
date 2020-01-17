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
package example

import (
	"context"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "example"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result  chan *diagnose.Result
	Message string `yaml:"message"`
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
func (d *Diagnostic) StartDiagnose(ctx context.Context,
	param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		// send a risk result
		obj := map[string]interface{}{
			"Mes": d.Message,
		}

		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelRisk,
			Title:    "example",
			ObjName:  "example-obj",
			ObjInfo:  obj,
			Desc:     d.Translator.Message("message", obj),
			Proposal: d.Translator.Message("proposal", nil),
		}

		// if any Error occur , send a failed result
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelFailed,
			Title:    "example",
			ObjName:  "example-obj",
			ObjInfo:  obj,
			Desc:     d.Translator.Message("message", obj),
			Proposal: d.Translator.Message("proposal", nil),
		}
	}()
	return d.result, nil
}
