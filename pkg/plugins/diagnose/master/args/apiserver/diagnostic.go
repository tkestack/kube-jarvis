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
package apiserver

import (
	"context"
	"fmt"
	"strconv"

	"tkestack.io/kube-jarvis/pkg/translate"

	"k8s.io/apimachinery/pkg/api/resource"

	v1 "k8s.io/api/core/v1"

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "kube-apiserver-args"
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
		for _, comp := range param.Resources.CoreComponents[cluster.ComponentApiserver] {
			d.checkOne(param.Resources, comp)
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) checkOne(resources *cluster.Resources, info cluster.Component) {
	if info.Error != nil {
		d.Logger.Errorf("check kube-apiserver on node %s get error : %v", info.Node, info.Error)
		return
	}

	if !info.IsRunning {
		d.Logger.Errorf("kube-apiserver on node %s not running ", info.Node)
		return
	}

	if info.Node == "" {
		d.Logger.Errorf("can not get node info for %s", info.Name)
		return
	}

	node := findNode(resources.Nodes, info.Node)
	mem := node.Status.Capacity.Memory()
	cpu := node.Status.Capacity.Cpu()
	riTarget := 400
	mriTarget := 200

	if mem.Cmp(resource.MustParse("14Gi")) > 1 && cpu.Cmp(resource.MustParse("6000m")) > 1 {
		riTarget = 800
		mriTarget = 300
	}

	if mem.Cmp(resource.MustParse("30Gi")) > 1 && cpu.Cmp(resource.MustParse("14000m")) > 1 {
		riTarget = 1500
		mriTarget = 500
	}

	if mem.Cmp(resource.MustParse("60Gi")) > 1 && cpu.Cmp(resource.MustParse("30000m")) > 1 {
		riTarget = 3000
		mriTarget = 1000
	}

	d.checkArgs(resources, info, "max-requests-inflight", 400, riTarget)
	d.checkArgs(resources, info, "max-mutating-requests-inflight", 200, mriTarget)
}

func (d *Diagnostic) checkArgs(resources *cluster.Resources, info cluster.Component, arg string, defVal, targetVal int) {
	nodeTotal := len(resources.Nodes.Items)
	obj := map[string]interface{}{
		"Name":      info.Name,
		"Node":      info.Node,
		"NodeTotal": nodeTotal,
		"Arg":       arg,
	}

	curVal := defVal
	curValStr := info.Args[arg]
	if curValStr != "" {
		curVal, _ = strconv.Atoi(curValStr)
	}

	obj["CurVal"] = curVal
	obj["TargetVal"] = targetVal

	level := diagnose.HealthyLevelGood
	desc := d.Translator.Message("good-desc", obj)
	proposal := translate.Message("")

	if curVal < targetVal {
		level = diagnose.HealthyLevelRisk
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

func findNode(node *v1.NodeList, nodeName string) *v1.Node {
	for _, n := range node.Items {
		if n.Name == nodeName {
			return &n
		}
	}
	return nil
}
