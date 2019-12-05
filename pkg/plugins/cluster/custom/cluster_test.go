package custom

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster/custom/compexplorer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
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

	cls.Components = map[string]*compexplorer.Auto{}
	if err := cls.Init(); err != nil {
		t.Fatalf(err.Error())
	}

	cls.compExps = map[string]compexplorer.Explorer{
		cluster.ComponentApiserver: &fakeComp{},
	}
	cls.nodeExecutor = &fakeNodeExecutor{success: true}

	if err := cls.SyncResources(); err != nil {
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

}
