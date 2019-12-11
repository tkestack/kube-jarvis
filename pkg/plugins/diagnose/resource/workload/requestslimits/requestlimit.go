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
package requestslimits

import (
	"context"
	"fmt"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v12 "k8s.io/api/core/v1"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "requests-limits"
)

// Diagnostic report the healthy of pods's resources requests limits configuration
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a requests-limits Diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		MetaData: meta,
		result:   make(chan *diagnose.Result, 1000),
	}
}

func (d *Diagnostic) Init() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) chan *diagnose.Result {
	d.param = &param
	go func() {
		defer diagnose.CommonDeafer(d.result)
		for _, pod := range d.param.Resources.Pods.Items {
			d.diagnosePod(pod)
		}
	}()
	return d.result
}

func (d *Diagnostic) diagnosePod(pod v12.Pod) {
	for _, c := range pod.Spec.Containers {
		if c.Resources.Limits.Memory().IsZero() ||
			c.Resources.Limits.Cpu().IsZero() ||
			c.Resources.Requests.Memory().IsZero() ||
			c.Resources.Requests.Cpu().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Title:    "Pods Requests Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Namespace, pod.Name),
				Desc:     d.Translator.Message("desc", nil),
				Proposal: d.Translator.Message("proposal", nil),
			}
			return
		}
	}
}
