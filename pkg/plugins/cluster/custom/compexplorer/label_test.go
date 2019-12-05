package compexplorer

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestDaemonSet_Component(t *testing.T) {
	fk := fake.NewSimpleClientset()
	pod1 := &v1.Pod{}
	pod1.Name = "p1"
	pod1.Namespace = "kube-system"
	pod1.Labels = map[string]string{
		"k8s-app": "p1",
	}

	pod1.Spec.NodeName = "node1"

	pod1.Spec.Containers = []v1.Container{
		{
			Name: "p1",
			Args: []string{
				"--a = 123",
				"--b = 321",
			},
		},
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

	l := NewLabelExp(logger.NewLogger(), fk, "kube-system", "p1", map[string]string{
		"k8s-app": "p1",
	})

	cmp, err := l.Component()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(cmp) != 1 {
		t.Fatalf("want 1 result")
	}

	if !cmp[0].IsRunning {
		t.Fatalf("wan Runing")
	}

	if cmp[0].Name != pod1.Name {
		t.Fatalf("want name %s ,but get %s", pod1.Name, cmp[0].Name)
	}

	if cmp[0].Node != pod1.Spec.NodeName {
		t.Fatalf("want nodeName %s, but get %s ", pod1.Spec.NodeName, cmp[0].Node)
	}

	if cmp[0].Args["a"] != "123" {
		t.Fatalf("want key a value 123 , but get %s", cmp[0].Args["a"])
	}

	if cmp[0].Args["b"] != "321" {
		t.Fatalf("want key b value 321 , but get %s", cmp[0].Args["a"])
	}
}
