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
package custom

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/compexplorer"
)

type fakeComp struct {
}

func (f *fakeComp) Component() ([]cluster.Component, error) {
	return []cluster.Component{
		{
			Name:      "kube-apiserver",
			IsRunning: true,
		},
	}, nil
}

func (f *fakeComp) Finish() error {
	return nil
}

type fakeNodeExecutor struct {
	success bool
}

func (f *fakeNodeExecutor) DoCmd(nodeName string, cmd []string) (string, string, error) {
	out := `kube-apiserver
-a=123
-b=321
`
	if !f.success {
		return "", "", nil
	}
	return out, "", nil
}

func (f *fakeNodeExecutor) Finish() error {
	return nil
}

func TestGetSysCtlMap(t *testing.T) {
	out := `
	a = 1
	b = 2 3 
`
	m := GetSysCtlMap(out)
	t.Logf("%+v", m)
	if len(m) != 2 {
		t.Fatalf("want 2 key")
	}

	if m["a"] != "1" {
		t.Fatalf("key a want value 1")
	}

	if m["b"] != "2 3" {
		t.Fatalf("key b want value 2")
	}
}

func TestCluster_Resources(t *testing.T) {
	fk := fake.NewSimpleClientset()
	pod := &v1.Pod{}
	pod.Name = "pod1"
	pod.Namespace = "kube-system"

	ns := &v1.Namespace{}
	ns.Name = "kube-system"
	if _, err := fk.CoreV1().Namespaces().Create(ns); err != nil {
		t.Fatalf(err.Error())
	}

	if _, err := fk.CoreV1().Pods(pod.Namespace).Create(pod); err != nil {
		t.Fatalf(err.Error())
	}

	node := &v1.Node{}
	node.Name = "node1"
	if _, err := fk.CoreV1().Nodes().Create(node); err != nil {
		t.Fatalf(err.Error())
	}

	cls := NewCluster(logger.NewLogger(), fk, nil).(*Cluster)
	if cls.CloudType() != Type {
		t.Fatalf("wrong cloud type")
	}

	if err := cls.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	cls.Components = map[string]*compexplorer.Auto{}
	cls.compExps = map[string]compexplorer.Explorer{
		cluster.ComponentApiserver: &fakeComp{},
	}
	cls.nodeExecutor = &fakeNodeExecutor{success: true}

	if err := cls.Init(context.Background(), plugins.NewProgress()); err != nil {
		t.Fatalf(err.Error())
	}

	res := cls.Resources()
	if len(res.Pods.Items) != 1 {
		t.Fatalf("want 1 Pods")
	}

	if len(res.Machines) != 1 {
		t.Fatalf("want 1 Machines")
	}

	if len(res.CoreComponents) != 1 {
		t.Fatalf("want 1 CoreComponents")
	}

	if err := cls.Finish(); err != nil {
		t.Fatalf(err.Error())
	}
}
