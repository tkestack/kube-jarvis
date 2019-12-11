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
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

// Bare get component information from cmd
type Bare struct {
	logger       logger.Logger
	cmdName      string
	nodes        []string
	nodeExecutor nodeexec.Executor
}

// NewBare create and int a StaticPods ComponentExecutor
func NewBare(logger logger.Logger, cmdName string, nodes []string, executor nodeexec.Executor) *Bare {
	return &Bare{
		logger:       logger,
		cmdName:      cmdName,
		nodes:        nodes,
		nodeExecutor: executor,
	}
}

// Component get cluster components
func (b *Bare) Component() ([]cluster.Component, error) {
	cmd := fmt.Sprintf("pgrep %s &&  cat /proc/`pgrep %s`/cmdline | xargs -0 | tr ' ' '\\n'", b.cmdName, b.cmdName)
	result := make([]cluster.Component, 0)
	lk := sync.Mutex{}
	conCtl := make(chan struct{}, 200)
	g := errgroup.Group{}

	for _, tempN := range b.nodes {
		n := tempN
		g.Go(func() error {
			conCtl <- struct{}{}
			defer func() { <-conCtl }()

			out, _, err := b.nodeExecutor.DoCmd(n, []string{
				"/bin/sh", "-c", cmd,
			})
			if err != nil {
				if !strings.Contains(err.Error(), "terminated with exit code") {
					b.logger.Errorf("do command on node %s failed :%v", n, err)
				}
				return err
			}

			cmp := cluster.Component{
				Name: b.cmdName,
				Node: n,
				Args: map[string]string{},
			}

			lines := strings.Split(out, "\n")
			for i, line := range lines {
				line = strings.TrimSpace(line)
				line = strings.TrimLeft(line, "-")
				if line == "" {
					continue
				}

				if i == 0 {
					cmp.IsRunning = true
					continue
				}

				spIndex := strings.IndexAny(line, "=")
				if spIndex == -1 {
					continue
				}

				k := line[0:spIndex]
				v := line[spIndex+1:]
				cmp.Args[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}

			lk.Lock()
			result = append(result, cmp)
			lk.Unlock()
			return nil
		})
	}
	_ = g.Wait()

	return result, nil
}
