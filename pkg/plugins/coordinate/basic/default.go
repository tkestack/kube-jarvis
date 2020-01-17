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
package basic

import (
	"context"
	"time"

	"tkestack.io/kube-jarvis/pkg/store"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

// Coordinator Coordinate diagnostics,exporters,evaluators with simple way
type Coordinator struct {
	cls         cluster.Cluster
	logger      logger.Logger
	diagnostics []diagnose.Diagnostic
	exporters   []export.Exporter
	progress    *plugins.Progress
	store       store.Store
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger,
	cls cluster.Cluster, st store.Store) coordinate.Coordinator {
	return &Coordinator{
		logger: logger,
		cls:    cls,
		store:  st,
	}
}

// Complete check and complete check config items
func (c *Coordinator) Complete() error {
	return nil
}

// AddDiagnostic add a diagnostic to Coordinator
func (c *Coordinator) AddDiagnostic(dia diagnose.Diagnostic) {
	c.diagnostics = append(c.diagnostics, dia)
}

// AddExporter add a Exporter to Coordinator
func (c *Coordinator) AddExporter(exporter export.Exporter) {
	c.exporters = append(c.exporters, exporter)
}

// Run will do all diagnostics, evaluations, then export it by exporters
func (c *Coordinator) Run(ctx context.Context) error {
	c.progress = plugins.NewProgress()
	c.progress.AddProgressUpdatedWatcher(func(p *plugins.Progress) {
		c.progress = p.Clone()
	})

	c.progress.CreateStep("diagnostic", "Diagnosing...", len(c.diagnostics))

	if err := c.cls.Init(ctx, c.progress); err != nil {
		return errors.Wrap(err, "init cluster failed")
	}

	c.logger.Infof("Start Diagnosing......")
	c.progress.SetCurStep("diagnostic")
	c.diagnostic(ctx)

	if err := c.cls.Finish(); err != nil {
		return errors.Wrapf(err, "finish cluster failed")
	}

	c.progress.Done()
	c.logger.Infof("Diagnosing done!")
	return nil
}

// Progress return the coordination progress
// if coordination has not start, an nil will be returned
func (c *Coordinator) Progress() *plugins.Progress {
	return c.progress
}

func (c *Coordinator) diagnostic(ctx context.Context) {
	result := export.NewAllResult()
	for _, dia := range c.diagnostics {
		resultChan, err := dia.StartDiagnose(ctx, diagnose.StartDiagnoseParam{
			CloudType: c.cls.CloudType(),
			Resources: c.cls.Resources(),
		})

		if err != nil {
			c.logger.Errorf("start diagnostic type[%s] name[%s] failed : %v",
				dia.Meta().Type, dia.Meta().Name, err)
			return
		}

		resultItem := export.NewDiagnosticResultItem(dia)
		for {
			s, ok := <-resultChan
			if !ok {
				break
			}
			resultItem.AddResult(s)
		}
		resultItem.EndTime = time.Now()

		result.AddDiagnosticResultItem(resultItem)
		c.progress.AddStepPercent("diagnostic", 1)
	}

	result.EndTime = time.Now()
	c.export(ctx, result)
}

func (c *Coordinator) export(ctx context.Context, r *export.AllResult) {
	g := errgroup.Group{}
	for _, tmp := range c.exporters {
		e := tmp
		g.Go(func() error {
			if err := e.Export(ctx, r); err != nil {
				c.logger.Errorf("%s export failed: %v", e.Meta().Name, err)
			}
			return nil
		})
	}
	_ = g.Wait()
}
