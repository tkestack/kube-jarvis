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
package etcd

import (
	"context"
	"fmt"
	"strconv"

	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "etcd-args"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
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
	go func() {
		defer diagnose.CommonDeafer(d.result)
		for _, comp := range param.Resources.CoreComponents[cluster.ComponentETCD] {
			d.checkQuota(param.Resources, comp)
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) checkQuota(resources *cluster.Resources, info cluster.Component) {
	arg := "quota-backend-bytes"
	curVal := 2 * 1024 * 1024 * 1024
	curValStr := info.Args[arg]
	if curValStr != "" {
		curVal, _ = strconv.Atoi(curValStr)
	}
	targetVal := 6 * 1024 * 1024 * 1024

	obj := map[string]interface{}{
		"Name":      info.Name,
		"Node":      info.Node,
		"Arg":       arg,
		"CurVal":    curVal,
		"TargetVal": targetVal,
	}

	level := diagnose.HealthyLevelGood
	desc := d.Translator.Message("good-desc", obj)
	proposal := translate.Message("")

	if curVal < targetVal {
		level = diagnose.HealthyLevelWarn
		desc = d.Translator.Message(fmt.Sprintf("%s-desc", arg), obj)
		proposal = d.Translator.Message(fmt.Sprintf("%s-proposal", arg), obj)
	}

	d.result <- &diagnose.Result{
		Level:    level,
		ObjName:  info.Name,
		ObjInfo:  obj,
		Title:    d.Translator.Message(fmt.Sprintf("%s-title", arg), nil),
		Desc:     desc,
		Proposal: proposal,
	}
}
