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
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
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
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger, cls cluster.Cluster) coordinate.Coordinator {
	return &Coordinator{
		logger:   logger,
		cls:      cls,
		progress: plugins.NewProgress(),
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
func (c *Coordinator) Run(ctx context.Context) {
	c.begin(ctx)
	c.progress.AddProgressUpdatedWatcher(func(p *plugins.Progress) {
		c.everyExporterDo(func(e export.Exporter) {
			c.logIfError(e.ProgressUpdated(ctx, p), "%s export progress update failed", e.Meta().Name)
		})
	})

	c.progress.CreateStep("diagnostic", "Diagnosing...", len(c.diagnostics))
	c.logIfError(c.cls.Init(ctx, c.progress), "init cluster failed")
	c.diagnostic(ctx)
	c.progress.Done()
	c.finish(ctx)
}

func (c *Coordinator) begin(ctx context.Context) {
	c.everyExporterDo(func(e export.Exporter) {
		c.logIfError(e.CoordinateBegin(ctx), "%s export coordinate begin", e.Meta().Name)
	})
}

func (c *Coordinator) finish(ctx context.Context) {
	c.everyExporterDo(func(e export.Exporter) {
		c.logIfError(e.CoordinateFinish(ctx), "%s export coordinate finish", e.Meta().Name)
	})
}

func (c *Coordinator) diagnostic(ctx context.Context) {
	for _, dia := range c.diagnostics {
		c.diagnosticBegin(ctx, dia)
		result, err := dia.StartDiagnose(ctx, diagnose.StartDiagnoseParam{
			CloudType: c.cls.CloudType(),
			Resources: c.cls.Resources(),
		})
		if err != nil {
			c.logger.Errorf("start diagnostic type[%s] name[%s] failed : %v", dia.Meta().Type, dia.Meta().Name, err)
			os.Exit(1)
		}

		for {
			s, ok := <-result
			if !ok {
				break
			}
			c.notifyDiagnosticResult(ctx, dia, s)
		}
		c.diagnosticFinish(ctx, dia)
		c.progress.AddStepPercent("diagnostic", 1)
	}
}

func (c *Coordinator) diagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) {
	c.everyExporterDo(func(e export.Exporter) {
		c.logIfError(e.DiagnosticBegin(ctx, dia), "%s export diagnose begin", e.Meta().Name)
	})
}

func (c *Coordinator) diagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) {
	c.everyExporterDo(func(e export.Exporter) {
		c.logIfError(e.DiagnosticFinish(ctx, dia), "%s export diagnose finish", e.Meta().Name)
	})
}

func (c *Coordinator) notifyDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) {
	c.everyExporterDo(func(e export.Exporter) {
		c.logIfError(e.DiagnosticResult(ctx, dia, result), "%s export diagnose result", e.Meta().Name)
	})
}

func (c *Coordinator) logIfError(err error, format string, args ...interface{}) {
	if err != nil {
		c.logger.Errorf("%s failed : %v", fmt.Sprintf(format, args...), err.Error())
	}
}

func (c *Coordinator) everyExporterDo(f func(e export.Exporter)) {
	g := errgroup.Group{}
	for _, tmp := range c.exporters {
		e := tmp
		g.Go(func() error {
			f(e)
			return nil
		})
	}
	_ = g.Wait()
}
