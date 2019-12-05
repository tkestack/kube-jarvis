package nodeexec

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestConfig_Executor(t *testing.T) {
	n := NewConfig()
	exe, err := n.Executor(logger.NewLogger(), fake.NewSimpleClientset(), nil)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, ok := exe.(*DaemonSetProxy)
	if !ok {
		t.Fatalf("should return an DaemonSetProxy Executor")
	}

	n.Type = "none"
	_, err = n.Executor(logger.NewLogger(), fake.NewSimpleClientset(), nil)
	if err != NoneExecutor {
		t.Fatalf("should get a UnKnowTypeErr")
	}
}
