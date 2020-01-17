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

	"github.com/pkg/errors"
	v12 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

// StaticPods get component information from static pod
type StaticPods struct {
	logger      logger.Logger
	cli         kubernetes.Interface
	podName     string
	namespace   string
	nodes       []string
	exec        nodeexec.Executor
	ExplorePods func(logger logger.Logger, name string,
		pods []v12.Pod, exec nodeexec.Executor) []cluster.Component
}

// NewStaticPods create and int a StaticPods ComponentExecutor
func NewStaticPods(logger logger.Logger, cli kubernetes.Interface, namespace string,
	podPrefix string, nodes []string, exe nodeexec.Executor) *StaticPods {
	return &StaticPods{
		logger:      logger,
		cli:         cli,
		podName:     podPrefix,
		namespace:   namespace,
		nodes:       nodes,
		exec:        exe,
		ExplorePods: ExplorePods,
	}
}

// Component get cluster components
func (s *StaticPods) Component() ([]cluster.Component, error) {
	result := make([]cluster.Component, 0)
	pods := make([]v12.Pod, 0)
	for _, n := range s.nodes {
		cmp := cluster.Component{
			Name: s.podName,
			Node: n,
		}

		podName := fmt.Sprintf("%s-%s", s.podName, n)
		pod, err := s.cli.CoreV1().Pods(s.namespace).Get(podName, v1.GetOptions{})
		if err != nil {
			if !k8serr.IsNotFound(err) {
				cmp.Error = errors.Wrapf(err, "get target pod %s failed", cmp.Name)
			}
			result = append(result, cmp)
			continue
		}

		pods = append(pods, *pod)
	}

	return append(result, s.ExplorePods(s.logger, s.podName, pods, s.exec)...), nil
}

// Finish will be called once every thing done
func (s *StaticPods) Finish() error {
	return nil
}
