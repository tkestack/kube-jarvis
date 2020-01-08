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
package pdb

import (
	"context"
	"fmt"

	"k8s.io/api/policy/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v12 "k8s.io/api/core/v1"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "pdb"
)

// Diagnostic report the healthy of pods's resources health check configuration
type Diagnostic struct {
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
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		uid2obj := make(map[types.UID]diagnose.MetaObject)
		outputs := make(map[types.UID]bool)
		podNum := make(map[types.UID]int32)
		for _, deploy := range d.param.Resources.Deployments.Items {
			deploy.Kind = "Deployment"
			uid2obj[deploy.UID] = deploy.DeepCopy()
			if deploy.Spec.Replicas == nil {
				podNum[deploy.UID] = 1
			} else {
				podNum[deploy.UID] = *deploy.Spec.Replicas
			}
		}
		for _, sts := range d.param.Resources.StatefulSets.Items {
			sts.Kind = "StatefulSet"
			uid2obj[sts.UID] = sts.DeepCopy()
			if sts.Spec.Replicas == nil {
				podNum[sts.UID] = 1
			} else {
				podNum[sts.UID] = *sts.Spec.Replicas
			}
		}
		for _, rs := range d.param.Resources.ReplicaSets.Items {
			rs.Kind = "ReplicaSet"
			uid2obj[rs.UID] = rs.DeepCopy()
			if rs.Spec.Replicas == nil {
				podNum[rs.UID] = 1
			} else {
				podNum[rs.UID] = *rs.Spec.Replicas
			}
		}
		for _, rc := range d.param.Resources.ReplicationControllers.Items {
			rc.Kind = "ReplicationController"
			uid2obj[rc.UID] = rc.DeepCopy()
			if rc.Spec.Replicas == nil {
				podNum[rc.UID] = 1
			} else {
				podNum[rc.UID] = *rc.Spec.Replicas
			}
		}
		for _, ds := range d.param.Resources.DaemonSets.Items {
			ds.Kind = "DaemonSet"
			uid2obj[ds.UID] = ds.DeepCopy()
		}

		for _, pod := range d.param.Resources.Pods.Items {
			pod.Kind = "Pod"
			rootOwner := diagnose.GetRootOwner(&pod, uid2obj)
			if rootOwner.GroupVersionKind().Kind == "DaemonSet" || rootOwner.GroupVersionKind().Kind == "Pod" {
				continue
			}

			if _, ok := outputs[rootOwner.GetUID()]; ok {
				continue
			}

			if podNum[rootOwner.GetUID()] <= 1 {
				continue
			}

			d.diagnosePod(&pod, rootOwner, d.param.Resources.PodDisruptionBudgets)
			outputs[rootOwner.GetUID()] = true
		}
	}()
	return d.result, nil
}

func getPodDisruptionBudgets(pod *v12.Pod, pdbList *v1beta1.PodDisruptionBudgetList) ([]v1beta1.PodDisruptionBudget, error) {
	if pod == nil || len(pod.Labels) == 0 {
		return nil, nil
	}

	var pdbs []v1beta1.PodDisruptionBudget
	for _, pdb := range pdbList.Items {
		if pdb.Namespace != pod.Namespace {
			continue
		}
		selector, err := v1.LabelSelectorAsSelector(pdb.Spec.Selector)
		if err != nil {
			continue
		}
		// If a PDB with a nil or empty selector creeps in, it should match nothing, not everything.
		if selector.Empty() || !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}

		pdbs = append(pdbs, pdb)
	}

	return pdbs, nil
}

func (d *Diagnostic) diagnosePod(pod *v12.Pod, rootOwner diagnose.MetaObject, pdbList *v1beta1.PodDisruptionBudgetList) {
	pdbs, err := getPodDisruptionBudgets(pod, pdbList)
	if err != nil {
		return
	}
	if len(pdbs) > 1 {
		// invalid
	} else if len(pdbs) == 0 {
		obj := map[string]interface{}{
			"Kind":      rootOwner.GroupVersionKind().Kind,
			"Namespace": rootOwner.GetNamespace(),
			"Name":      rootOwner.GetName(),
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
}
