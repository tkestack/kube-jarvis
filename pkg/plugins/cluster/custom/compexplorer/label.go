package compexplorer

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// LabelExp get select pod with labels and try get component from
type LabelExp struct {
	logger    logger.Logger
	cli       kubernetes.Interface
	namespace string
	name      string
	labels    map[string]string
}

// NewLabelExp create and init a LabelExp Component
func NewLabelExp(logger logger.Logger, cli kubernetes.Interface, namespace string, cmpName string, labels map[string]string) *LabelExp {
	return &LabelExp{
		logger:    logger,
		cli:       cli,
		name:      cmpName,
		labels:    labels,
		namespace: namespace,
	}
}

// Component get cluster components
func (n *LabelExp) Component() ([]cluster.Component, error) {
	result := make([]cluster.Component, 0)
	pods, err := n.cli.CoreV1().Pods(n.namespace).List(v12.ListOptions{
		LabelSelector: labels.FormatLabels(n.labels),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "get pods failed")
	}

	for _, pod := range pods.Items {
		result = append(result, cluster.Component{
			Name:      pod.Name,
			Node:      pod.Spec.NodeName,
			Args:      n.getArgs(&pod),
			IsRunning: true,
		})
	}

	return result, nil
}

func (n *LabelExp) getArgs(pod *v1.Pod) map[string]string {
	result := make(map[string]string)
	for _, c := range pod.Spec.Containers {
		if c.Name == n.name {
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
