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
package compexplorer

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

func TestDaemonSet_Component(t *testing.T) {
	fk := fake.NewSimpleClientset()
	pod1 := &v1.Pod{}
	pod1.Name = "p1"
	pod1.Namespace = "kube-system"
	pod1.Labels = map[string]string{
		"k8s-app": "p1",
	}

	if _, err := fk.CoreV1().Pods("kube-system").Create(pod1); err != nil {
		t.Fatalf(err.Error())
	}

	pod2 := &v1.Pod{}
	pod2.Namespace = "kube-system"
	pod2.Name = "p2"
	pod2.Labels = map[string]string{
		"k8s-app": "p2",
	}

	if _, err := fk.CoreV1().Pods("kube-system").Create(pod2); err != nil {
		t.Fatalf(err.Error())
	}

	l := NewLabelExp(logger.NewLogger(), fk, "kube-system", "p1", nil, nil)
	l.explorePods = func(logger logger.Logger, name string, pods []v1.Pod, exec nodeexec.Executor) []cluster.Component {
		if name != "p1" {
			t.Fatalf("name want p1 but get %s", name)
		}

		if len(pods) != 1 {
			t.Fatalf("want 1 pods but get %d", len(pods))
		}

		if pods[0].Name != "p1" {
			t.Fatalf("want p1 pod but get %s", pods[0].Name)
		}

		return nil
	}
}
