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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

func TestStaticPods_Component(t *testing.T) {
	fk := fake.NewSimpleClientset()
	total := 2
	nodes := make([]string, 0)
	for i := 0; i < total; i++ {
		n := &v1.Node{}
		n.Name = fmt.Sprintf("10.0.0.%d", i)
		if _, err := fk.CoreV1().Nodes().Create(n); err != nil {
			t.Fatalf("create master %s failed", n.Name)
		}

		nodes = append(nodes, n.Name)

		if i == 0 {
			continue
		}

		pod := &v1.Pod{}
		pod.Spec.NodeName = n.Name
		pod.Namespace = "kube-system"
		pod.Name = fmt.Sprintf("test-%s", n.Name)

		if _, err := fk.CoreV1().Pods("kube-system").Create(pod); err != nil {
			t.Fatalf(err.Error())
		}
	}

	sd := NewStaticPods(logger.NewLogger(), fk, "kube-system", "test", nodes, nil)
	sd.ExplorePods = func(logger logger.Logger, name string, pods []v1.Pod, exec nodeexec.Executor) []cluster.Component {
		if name != "test" {
			t.Fatalf("want name test but get %s", name)
		}

		if len(pods) != 1 {
			t.Fatalf("want 1 pods but get %d", len(pods))
		}

		return nil
	}

	cmp, err := sd.Component()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(cmp) != 1 {
		t.Fatalf("want 2 result but get %d", len(cmp))
	}

	if cmp[0].IsRunning {
		t.Fatalf("want cmponent not running ")
	}
}
