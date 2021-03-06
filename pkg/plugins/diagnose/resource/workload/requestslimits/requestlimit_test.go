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
	"testing"

	appv1 "k8s.io/api/apps/v1"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	res.Deployments = &appv1.DeploymentList{}
	res.ReplicaSets = &appv1.ReplicaSetList{}
	res.ReplicationControllers = &v1.ReplicationControllerList{}
	res.StatefulSets = &appv1.StatefulSetList{}
	res.DaemonSets = &appv1.DaemonSetList{}
	res.Pods = &v1.PodList{}

	kubeDNSlimit := make(map[v1.ResourceName]resource.Quantity)
	kubeDNSRequest := make(map[v1.ResourceName]resource.Quantity)

	kubeDNSlimit[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSlimit[v1.ResourceMemory] = resource.MustParse("170M")
	kubeDNSRequest[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSRequest[v1.ResourceMemory] = resource.MustParse("30M")

	pod := v1.Pod{}
	pod.Name = "pod1"
	pod.Namespace = "default"
	pod.UID = "pod1-uid"
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
			Resources: v1.ResourceRequirements{
				Limits:   kubeDNSlimit,
				Requests: kubeDNSRequest,
			},
		},
	}
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
