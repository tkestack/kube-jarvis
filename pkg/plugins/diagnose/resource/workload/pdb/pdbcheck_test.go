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
	"testing"

	appv1 "k8s.io/api/apps/v1"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
)

func TestPDBCheckDiagnostic_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	dep := appv1.Deployment{}
	dep.Name = "kube-dns"
	rs := int32(2)
	dep.Spec.Replicas = &rs
	dep.UID = "1"

	dep2 := appv1.Deployment{}
	dep2.Name = "kube-dns2"
	dep2.Spec.Replicas = &rs
	dep2.UID = "2"

	res.Deployments = &appv1.DeploymentList{
		Items: []appv1.Deployment{
			dep,
			dep2,
		},
	}
	res.ReplicaSets = &appv1.ReplicaSetList{}
	res.ReplicationControllers = &v1.ReplicationControllerList{}
	res.StatefulSets = &appv1.StatefulSetList{}
	res.DaemonSets = &appv1.DaemonSetList{}

	res.Pods = &v1.PodList{}
	res.PodDisruptionBudgets = &v1beta1.PodDisruptionBudgetList{}

	pdb := v1beta1.PodDisruptionBudget{}
	pdb.Name = "test-pdb"
	pdb.Namespace = "default"
	minAvailable := intstr.FromInt(1)
	pdb.Spec.MinAvailable = &minAvailable
	pdb.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{"test-pdb-key": "test-pdb-value"},
	}
	res.PodDisruptionBudgets.Items = append(res.PodDisruptionBudgets.Items, pdb)

	pod := v1.Pod{}
	pod.Name = "pod1"
	pod.Namespace = "default"
	pod.UID = "pod1-uid"
	pod.Labels = map[string]string{"test-pdb-key": "test-pdb-value"}
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
		},
	}

	owner := metav1.OwnerReference{}
	owner.Name = dep.Name
	owner.Kind = "Deployment"
	ctl := true
	owner.Controller = &ctl
	owner.UID = "1"

	pod.OwnerReferences = []metav1.OwnerReference{owner}
	res.Pods.Items = append(res.Pods.Items, pod)

	pod = v1.Pod{}
	pod.Name = "pod2"
	pod.Namespace = "default"
	pod.UID = "pod2-uid"
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
		},
	}
	owner2 := metav1.OwnerReference{}
	owner2.Name = dep2.Name
	owner2.Kind = "Deployment"
	owner2.Controller = &ctl
	owner2.UID = "2"
	pod.OwnerReferences = []metav1.OwnerReference{owner2}
	res.Pods.Items = append(res.Pods.Items, pod)

	d := NewDiagnostic(&diagnose.MetaData{
		MetaData: plugins.MetaData{
			Translator: translate.NewFake(),
		},
	})

	if err := d.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	result, _ := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
		Resources: res,
	})

	total := 0
	for {
		s, ok := <-result
		if !ok {
			break
		}
		total++

		t.Logf("%+v", *s)
	}
	if total != 1 {
		t.Fatalf("should return 1 result")
	}
}
