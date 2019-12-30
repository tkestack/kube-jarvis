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
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
)

func TestExplorePods(t *testing.T) {
	cases := []struct {
		useNodeExp bool
		cmpRunning bool
	}{
		{
			cmpRunning: false,
		},
		{
			cmpRunning: true,
			useNodeExp: true,
		}, {
			cmpRunning: true,
			useNodeExp: false,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			pod := v1.Pod{}
			pod.Name = "kube-apiserver"
			pod.Spec.NodeName = "10.0.0.1"
			pod.Spec.Containers = []v1.Container{
				{
					Name: "kube-apiserver",
					Args: []string{
						"-a=c1",
						"-b=c2",
					},
				},
			}

			if cs.cmpRunning {
				pod.Status = v1.PodStatus{
					Conditions: []v1.PodCondition{
						{
							Type:   v1.PodReady,
							Status: v1.ConditionTrue,
						},
					},
				}
			}

			results := ExplorePods(logger.NewLogger(), "kube-apiserver", []v1.Pod{pod}, &fakeNodeExecutor{success: true})
			if !cs.useNodeExp {
				results = ExplorePods(logger.NewLogger(), "kube-apiserver", []v1.Pod{pod}, nil)
			}

			if len(results) != 1 {
				t.Fatalf("want 1 result but get %d", len(results))
			}

			cmp := results[0]
			if cmp.Pod == nil {
				t.Fatalf("component pod should not nil")
			}

			if !cs.cmpRunning {
				if cmp.IsRunning {
					t.Fatalf("component should not running")
				}
				return
			}

			if len(cmp.Args) != 2 {
				t.Fatalf("want 2 args")
			}

			if cs.useNodeExp {
				if cmp.Args["a"] != "123" {
					t.Fatalf("key a want value 123 but get %s", cmp.Args["a"])
				}

				if cmp.Args["b"] != "321" {
					t.Fatalf("key a want value 123 but get %s", cmp.Args["b"])
				}
			} else {
				if cmp.Args["a"] != "c1" {
					t.Fatalf("key a want value c1 but get %s", cmp.Args["a"])
				}

				if cmp.Args["b"] != "c2" {
					t.Fatalf("key b want value c2 but get %s", cmp.Args["b"])
				}
			}
		})
	}
}

func TestGetPodArgs(t *testing.T) {
	pod := &v1.Pod{}
	pod.Name = "test"
	pod.Spec.Containers = []v1.Container{
		{
			Name: "test",
			Args: []string{
				"-a=123",
				"-b=321",
			},
		},
	}

	m := GetPodArgs("test", pod)
	if len(m) != 2 {
		t.Fatalf("should return 2 args")
	}

	if m["a"] != "123" {
		t.Fatalf("key a want 123 ,but get %s", m["a"])
	}

	if m["b"] != "321" {
		t.Fatalf("key a want 123 ,but get %s", m["a"])
	}
}
