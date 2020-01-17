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
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"tkestack.io/kube-jarvis/pkg/logger"
)

func TestAuto_Init(t *testing.T) {
	cases := []struct {
		masterNodes bool
		nodes       []string
		labels      map[string]string
	}{
		{
			masterNodes: true,
			nodes:       []string{},
			labels:      map[string]string{},
		},
		{
			masterNodes: false,
			nodes:       []string{},
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs), func(t *testing.T) {
			fk := fake.NewSimpleClientset()
			for i := 0; i < 6; i++ {
				node := &v1.Node{}
				node.Name = fmt.Sprintf("10.0.0.%d", i)

				if i < 3 {
					node.Labels = map[string]string{
						"node-role.kubernetes.io/master": "true",
					}
				}

				if _, err := fk.CoreV1().Nodes().Create(node); err != nil {
					t.Fatalf(err.Error())
				}
			}

			a := NewAuto("kube-apiserver", cs.masterNodes)
			a.Nodes = cs.nodes
			if len(cs.labels) != 0 {
				a.Labels = cs.labels
			}
			if err := a.Init(logger.NewLogger(), fk, &fakeNodeExecutor{success: true}); err != nil {
				t.Fatalf(err.Error())
			}

			if cs.masterNodes && len(a.Nodes) != 3 {
				t.Fatalf("want 3 nodes")
			}

			if !cs.masterNodes && len(cs.nodes) == 0 && len(a.Nodes) != 6 {
				t.Fatalf("want 6 nodes")
			}
		})
	}
}

func TestAuto_Component(t *testing.T) {
	fk := fake.NewSimpleClientset()
	a := NewAuto("kube-apiserver", false)
	a.Nodes = []string{"node1"}
	if err := a.Init(logger.NewLogger(), fk, &fakeNodeExecutor{success: true}); err != nil {
		t.Fatalf(err.Error())
	}

	cmp, err := a.Component()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(cmp) != 1 {
		t.Fatalf("want len 1 but get %d", len(cmp))
	}

	if !cmp[0].IsRunning {
		t.Fatalf("IsRuning wrong")
	}

	if cmp[0].Args["a"] != "123" {
		t.Fatalf("want key a valuer 123 but get %s", cmp[0].Args["a"])
	}

	if cmp[0].Args["b"] != "321" {
		t.Fatalf("want key a valuer 321 but get %s", cmp[0].Args["a"])
	}
}
