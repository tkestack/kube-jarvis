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

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"k8s.io/apimachinery/pkg/types"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "workload-ha"
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
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		uid2obj := make(map[types.UID]diagnose.MetaObject)
		for _, deploy := range d.param.Resources.Deployments.Items {
			deploy.Kind = "Deployment"
			if deploy.Spec.Replicas == nil || *deploy.Spec.Replicas == 1 {
				continue
			}

			if d.Filter.Filtered(deploy.Name, "Deployment", deploy.Name) {
				continue
			}
			uid2obj[deploy.UID] = deploy.DeepCopy()
		}

		for _, sts := range d.param.Resources.StatefulSets.Items {
			sts.Kind = "StatefulSet"
			if sts.Spec.Replicas == nil || *sts.Spec.Replicas == 1 {
				continue
			}

			if d.Filter.Filtered(sts.Namespace, "StatefulSet", sts.Name) {
				continue
			}

			uid2obj[sts.UID] = sts.DeepCopy()
		}

		for _, rs := range d.param.Resources.ReplicaSets.Items {
			rs.Kind = "ReplicaSet"
			if rs.Spec.Replicas == nil || *rs.Spec.Replicas == 1 {
				continue
			}

			if d.Filter.Filtered(rs.Namespace, "ReplicaSet", rs.Name) {
				continue
			}

			uid2obj[rs.UID] = rs.DeepCopy()
		}

		for _, rc := range d.param.Resources.ReplicationControllers.Items {
			rc.Kind = "ReplicationController"
			if rc.Spec.Replicas == nil || *rc.Spec.Replicas == 1 {
				continue
			}

			if d.Filter.Filtered(rc.Namespace, "ReplicationController", rc.Name) {
				continue
			}

			uid2obj[rc.UID] = rc.DeepCopy()
		}

		lastNode := map[types.UID]string{}
		for _, pod := range d.param.Resources.Pods.Items {
			pod.Kind = "Pod"
			rootOwner := diagnose.GetRootOwner(&pod, uid2obj)
			if rootOwner.GroupVersionKind().Kind == "Pod" {
				continue
			}

			last := lastNode[rootOwner.GetUID()]
			if last == "pass" {
				continue
			}

			if last == "" {
				lastNode[rootOwner.GetUID()] = pod.Spec.NodeName
			} else if last != pod.Spec.NodeName {
				lastNode[rootOwner.GetUID()] = "pass"
			}
		}

		for uid, node := range lastNode {
			if node == "pass" {
				continue
			}
			d.diagnoseBad(uid2obj[uid], node)
		}

	}()
	return d.result, nil
}

func (d *Diagnostic) diagnoseBad(rootOwner diagnose.MetaObject, node string) {
	obj := map[string]interface{}{
		"Kind":      rootOwner.GroupVersionKind().Kind,
		"Namespace": rootOwner.GetNamespace(),
		"Name":      rootOwner.GetName(),
		"Node":      node,
	}

	d.result <- &diagnose.Result{
		Level:    diagnose.HealthyLevelWarn,
		ObjName:  fmt.Sprintf("%s:%s", rootOwner.GetNamespace(), rootOwner.GetName()),
		ObjInfo:  obj,
		Title:    d.Translator.Message("title", nil),
		Desc:     d.Translator.Message("desc", obj),
		Proposal: d.Translator.Message("proposal", obj),
	}
}
