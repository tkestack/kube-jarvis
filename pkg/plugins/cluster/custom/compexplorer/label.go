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
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"strings"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
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
