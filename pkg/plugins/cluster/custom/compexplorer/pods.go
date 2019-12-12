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
	"golang.org/x/sync/errgroup"
	v1 "k8s.io/api/core/v1"
	"strings"
	"sync"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

// ExplorePod explore a component from k8s pods
func ExplorePods(logger logger.Logger, name string, pods []v1.Pod, exec nodeexec.Executor) []cluster.Component {
	result := make([]cluster.Component, 0)
	lk := sync.Mutex{}
	g := errgroup.Group{}
	conCtl := make(chan struct{}, 200)

	for _, tempPod := range pods {
		pod := tempPod
		g.Go(func() error {
			conCtl <- struct{}{}
			defer func() { <-conCtl }()

			r := cluster.Component{
				Name: pod.Name,
				Node: pod.Spec.NodeName,
				Args: GetPodArgs(name, &pod),
				Pod:  &pod,
			}

			// see Ready as component Running
			for _, c := range pod.Status.Conditions {
				if c.Type == v1.PodReady && c.Status == v1.ConditionTrue {
					r.IsRunning = true
					break
				}
			}

			// we try to get args via node executor
			if exec != nil {
				bare := NewBare(logger, name, []string{pod.Spec.NodeName}, exec)
				cmp, err := bare.Component()
				if err != nil {
					logger.Errorf("try get component detail via node executor failed : %v", err)
				} else {
					if len(cmp) == 0 {
						logger.Errorf("can not found target component %s on node %s via node executor", name, pod.Spec.NodeName)
					} else if len(cmp[0].Args) == 0 {
						logger.Errorf("found target component %s on node %s ,but get empty args", name, pod.Spec.NodeName)
					} else {
						r.Args = cmp[0].Args
					}
				}
			}

			lk.Lock()
			result = append(result, r)
			lk.Unlock()
			return nil
		})
	}

	_ = g.Wait()

	return result
}

// GetPodArgs try get args from pod
func GetPodArgs(name string, pod *v1.Pod) map[string]string {
	result := make(map[string]string)
	for _, c := range pod.Spec.Containers {
		if c.Name == name {
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
