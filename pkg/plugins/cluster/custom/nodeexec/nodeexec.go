package nodeexec

import (
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

var (
	UnKnowTypeErr = fmt.Errorf("unknow node executor type")
	NoneExecutor  = fmt.Errorf("none executor")
)

// Executor get machine information
type Executor interface {
	// DoCmd do cmd on node and return output
	DoCmd(nodeName string, cmd []string) (string, string, error)
}

// Config is the config of node executor
type Config struct {
	Type      string
	Namespace string
	DaemonSet string
}

func NewConfig() *Config {
	return &Config{
		Type:      "proxy",
		Namespace: "kube-jarvis",
		DaemonSet: "kube-jarvis-agent",
	}
}

func (n *Config) Executor(logger logger.Logger, cli kubernetes.Interface, config *restclient.Config) (Executor, error) {
	switch n.Type {
	case "proxy":
		return NewDaemonSetProxy(logger, cli, config, n.Namespace, n.DaemonSet)
	case "none":
		return nil, NoneExecutor
	}
	return nil, UnKnowTypeErr
}
