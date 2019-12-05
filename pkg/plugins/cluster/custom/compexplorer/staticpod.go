package compexplorer

import (
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/pkg/errors"
	v12 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// StaticPods get component information from static pod
type StaticPods struct {
	logger    logger.Logger
	cli       kubernetes.Interface
	podName   string
	namespace string
	nodes     []string
}

// NewStaticPods create and int a StaticPods ComponentExecutor
func NewStaticPods(logger logger.Logger, cli kubernetes.Interface, namespace string, podPrefix string, nodes []string) *StaticPods {
	return &StaticPods{
		logger:    logger,
		cli:       cli,
		podName:   podPrefix,
		namespace: namespace,
		nodes:     nodes,
	}
}

// Component get cluster components
func (s *StaticPods) Component() ([]cluster.Component, error) {
	result := make([]cluster.Component, 0)
	for _, n := range s.nodes {
		cmp := cluster.Component{
			Name: s.podName,
			Node: n,
		}

		pod, err := s.cli.CoreV1().Pods(s.namespace).Get(fmt.Sprintf("%s-%s", s.podName, n), v1.GetOptions{})
		if err != nil {
			if k8serr.IsNotFound(err) {
				result = append(result, cmp)
				continue
			}
			return nil, errors.Wrapf(err, "get target pod %s failed", cmp.Name)
		}

		cmp.IsRunning = true
		cmp.Args = s.getArgs(pod)
		result = append(result, cmp)
	}
	return result, nil
}

func (s *StaticPods) getArgs(pod *v12.Pod) map[string]string {
	result := make(map[string]string)
	for _, c := range pod.Spec.Containers {
		if c.Name == s.podName {
			for _, arg := range c.Args {
				arg = strings.TrimLeft(arg, "-")
				spIndex := strings.IndexAny(arg, "=")
				if spIndex == -1 {
					continue
				}

				k := arg[0:spIndex]
				v := arg[spIndex+1:]
				result[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}
		}
	}
	return result
}
