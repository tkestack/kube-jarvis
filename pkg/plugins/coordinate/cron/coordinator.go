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
package cron

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"sync"
	"time"
	"tkestack.io/kube-jarvis/pkg/httpserver"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate/basic"
)

const (
	StateFailed  = "failed"
	StateRunning = "running"
	StatePending = "pending"
)

// Coordinator Coordinate diagnostics,exporters,evaluators with simple way
type Coordinator struct {
	Cron    string
	WalPath string
	coordinate.Coordinator
	state    string
	runLock  sync.Mutex
	cronCtl  *cron.Cron
	cronLock sync.Mutex
	waitRun  chan struct{}
	logger   logger.Logger
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger, cls cluster.Cluster) coordinate.Coordinator {
	c := &Coordinator{
		Coordinator: basic.NewCoordinator(logger, cls),
		waitRun:     make(chan struct{}),
		logger:      logger,
	}
	httpserver.HandleFunc("/coordinator/cron/run", c.runOnceHandler)
	httpserver.HandleFunc("/coordinator/cron/period", c.periodHandler)
	httpserver.HandleFunc("/coordinator/cron/state", c.stateHandler)
	return c
}

// Complete check and complete config items
func (c *Coordinator) Complete() error {
	return c.Coordinator.Complete()
}

// Run will do all diagnostics, evaluations, then export it by exporters
func (c *Coordinator) Run(ctx context.Context) error {
	if c.Cron != "" {
		c.cronCtl = cron.New()
		_, _ = c.cronCtl.AddFunc(c.Cron, c.cronDo)
		c.cronCtl.Start()
	}

	// check for wal file to auto start once we start
	// this is to ensure that the program automatically retries when it restarts
	go func() {
		if c.WalPath == "" {
			return
		}
		_, err := os.Stat(c.walFile()) //os.Stat获取文件信息
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			c.logger.Errorf("state wal file failed: %v", err)
		}
		c.logger.Infof("wal file exist, auto retry")
		c.tryStartRun()
	}()

	// start waiting for run
	for {
		select {
		case <-ctx.Done():
			c.logger.Infof("context done,coordinator exited")
			return ctx.Err()
		case <-c.waitRun:
		}
		c.runStart()
		if err := c.Coordinator.Run(ctx); err != nil {
			c.logger.Errorf("run failed: %v", err)
			c.runDone(false)
		} else {
			c.runDone(true)
		}
	}
}

func (c *Coordinator) runStart() {
	_, _ = os.Create(c.walFile())
}

func (c *Coordinator) runDone(success bool) {
	if c.WalPath != "" {
		_ = os.Remove(c.walFile())
	}
	c.runLock.Lock()
	defer c.runLock.Unlock()
	if success {
		c.state = StatePending
	} else {
		c.state = StateFailed
	}
}

func (c *Coordinator) walFile() string {
	return fmt.Sprintf("%s/kube-jarvis.wal", c.WalPath)
}

func (c *Coordinator) tryStartRun() bool {
	c.runLock.Lock()
	defer c.runLock.Unlock()

	if c.state == StateRunning {
		return false
	}

	c.state = StateRunning
	c.waitRun <- struct{}{}
	return true
}

func (c *Coordinator) cronDo() {
	for {
		if c.tryStartRun() {
			break
		}
		time.Sleep(time.Second * 1)
	}
}
