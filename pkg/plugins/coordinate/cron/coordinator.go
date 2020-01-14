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
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/robfig/cron/v3"
	"tkestack.io/kube-jarvis/pkg/httpserver"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate/basic"
	"tkestack.io/kube-jarvis/pkg/store"
)

const (
	StateFailed  = "failed"
	StateRunning = "running"
	StatePending = "pending"
)

// Coordinator Coordinate diagnostics,exporters,evaluators with simple way
type Coordinator struct {
	Cron string

	coordinate.Coordinator
	state    string
	runLock  sync.Mutex
	cronCtl  *cron.Cron
	cronLock sync.Mutex
	waitRun  chan struct{}
	logger   logger.Logger
	store    store.Store
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger, cls cluster.Cluster, st store.Store) coordinate.Coordinator {
	c := &Coordinator{
		Coordinator: basic.NewCoordinator(logger, cls, st),
		waitRun:     make(chan struct{}),
		logger:      logger,
		state:       StatePending,
		store:       st,
	}
	return c
}

// Complete check and complete config items
func (c *Coordinator) Complete() error {
	httpserver.HandleFunc("/coordinator/cron/run", c.runOnceHandler)
	httpserver.HandleFunc("/coordinator/cron/period", c.periodHandler)
	httpserver.HandleFunc("/coordinator/cron/state", c.stateHandler)
	if _, err := c.store.CreateSpace("cron"); err != nil {
		return errors.Wrap(err, "create store space failed")
	}
	return c.Coordinator.Complete()
}

// Run will do all diagnostics, evaluations, then export it by exporters
func (c *Coordinator) Run(ctx context.Context) error {
	if c.Cron != "" {
		c.cronCtl = cron.New(cron.WithSeconds())
		_, _ = c.cronCtl.AddFunc(c.Cron, c.cronDo)
		c.cronCtl.Start()
	}

	// check for auto start once we start
	// this is to ensure that the program automatically retries when it restarts
	go c.tryAutoStart()

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

func (c *Coordinator) tryAutoStart() {
	v, exist, err := c.store.Get("cron", "state")
	if err != nil {
		c.logger.Errorf("try auto start failed: %v", err.Error())
		return
	}

	if !exist || v != StateRunning {
		return
	}
	c.tryStartRun()
}

func (c *Coordinator) runStart() {
	_ = c.store.Set("cron", "state", StateRunning)
}

func (c *Coordinator) runDone(success bool) {
	_ = c.store.Set("cron", "state", StatePending)
	c.runLock.Lock()
	defer c.runLock.Unlock()
	if success {
		c.state = StatePending
	} else {
		c.state = StateFailed
	}
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
