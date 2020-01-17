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
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

// LabelExp get select pod with labels and try get component from
type LabelExp struct {
	logger      logger.Logger
	cli         kubernetes.Interface
	namespace   string
	name        string
	labels      map[string]string
	exec        nodeexec.Executor
	explorePods func(logger logger.Logger, name string,
		pods []v1.Pod, exec nodeexec.Executor) []cluster.Component
}

// NewLabelExp create and init a LabelExp Component
func NewLabelExp(logger logger.Logger, cli kubernetes.Interface,
	namespace string, cmpName string,
	labels map[string]string, exec nodeexec.Executor) *LabelExp {
	if len(labels) == 0 {
		labels = map[string]string{
			"k8s-app": cmpName,
		}
	}

	return &LabelExp{
		logger:      logger,
		cli:         cli,
		name:        cmpName,
		labels:      labels,
		namespace:   namespace,
		exec:        exec,
		explorePods: ExplorePods,
	}
}

// Component get cluster components
func (l *LabelExp) Component() ([]cluster.Component, error) {
	pods, err := l.cli.CoreV1().Pods(l.namespace).List(v12.ListOptions{
		LabelSelector: labels.FormatLabels(l.labels),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "get pods failed")
	}

	return l.explorePods(l.logger, l.name, pods.Items, l.exec), nil
}

// Finish will be called once every thing done
func (l *LabelExp) Finish() error {
	return nil
}
