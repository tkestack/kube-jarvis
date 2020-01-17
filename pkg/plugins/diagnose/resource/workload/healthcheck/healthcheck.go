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
package healthcheck

import (
	"context"
	"fmt"

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"k8s.io/apimachinery/pkg/types"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v12 "k8s.io/api/core/v1"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "health-check"
)

// Diagnostic report the healthy of pods's resources health check configuration
type Diagnostic struct {
	Filter cluster.ResourcesFilter
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a health check Diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		MetaData: meta,
		result:   make(chan *diagnose.Result, 1000),
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return d.Filter.Compile()
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context,
	param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		uid2obj := make(map[types.UID]diagnose.MetaObject)
		outputs := make(map[types.UID]bool)
		for _, deploy := range d.param.Resources.Deployments.Items {
			deploy.Kind = "Deployment"
			uid2obj[deploy.UID] = deploy.DeepCopy()
		}
		for _, sts := range d.param.Resources.StatefulSets.Items {
			sts.Kind = "StatefulSet"
			uid2obj[sts.UID] = sts.DeepCopy()
		}
		for _, rs := range d.param.Resources.ReplicaSets.Items {
			rs.Kind = "ReplicaSet"
			uid2obj[rs.UID] = rs.DeepCopy()
		}
		for _, rc := range d.param.Resources.ReplicationControllers.Items {
			rc.Kind = "ReplicationController"
			uid2obj[rc.UID] = rc.DeepCopy()
		}
		for _, ds := range d.param.Resources.DaemonSets.Items {
			ds.Kind = "DaemonSet"
			uid2obj[ds.UID] = ds.DeepCopy()
		}

		for _, pod := range d.param.Resources.Pods.Items {
			pod.Kind = "Pod"
			rootOwner := diagnose.GetRootOwner(&pod, uid2obj)
			if _, ok := outputs[rootOwner.GetUID()]; ok {
				continue
			}

			if d.Filter.Filtered(rootOwner.GetNamespace(),
				rootOwner.GroupVersionKind().Kind, rootOwner.GetName()) {
				continue
			}

			d.diagnosePod(&pod, rootOwner)
			outputs[rootOwner.GetUID()] = true
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) diagnosePod(pod *v12.Pod, rootOwner diagnose.MetaObject) {
	for _, c := range pod.Spec.Containers {
		if c.ReadinessProbe == nil || c.LivenessProbe == nil {
			obj := map[string]interface{}{
				"Kind":      rootOwner.GroupVersionKind().Kind,
				"Namespace": rootOwner.GetNamespace(),
				"Name":      rootOwner.GetName(),
			}

			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelRisk,
				ObjName:  fmt.Sprintf("%s:%s", rootOwner.GetNamespace(), rootOwner.GetName()),
				ObjInfo:  obj,
				Title:    d.Translator.Message("title", nil),
				Desc:     d.Translator.Message("desc", obj),
				Proposal: d.Translator.Message("proposal", obj),
			}
			return
		}
	}
}
