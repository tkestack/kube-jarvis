package basic

import (
	"context"
	"fmt"
	"os"
	"tkestack.io/kube-jarvis/pkg/logger"
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
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger, cls cluster.Cluster) coordinate.Coordinator {
	return &Coordinator{
		logger: logger,
		cls:    cls,
	}
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
	if err := c.cls.SyncResources(); err != nil {
		c.logger.Errorf("fetch resources failed :%v", err)
		os.Exit(1)
	}
	c.begin(ctx)
	c.diagnostic(ctx)
	c.finish(ctx)
}

func (c *Coordinator) begin(ctx context.Context) {
	for _, e := range c.exporters {
		c.logIfError(e.CoordinateBegin(ctx), "%s export coordinate begin", e.Meta().Name)
	}
}

func (c *Coordinator) finish(ctx context.Context) {
	for _, e := range c.exporters {
		c.logIfError(e.CoordinateFinish(ctx), "%s export coordinate finish", e.Meta().Name)
	}
}

func (c *Coordinator) diagnostic(ctx context.Context) {
	for _, dia := range c.diagnostics {
		c.diagnosticBegin(ctx, dia)
		result := dia.StartDiagnose(ctx, diagnose.StartDiagnoseParam{
			CloudType: c.cls.CloudType(),
			Resources: c.cls.Resources(),
		})

		for {
			s, ok := <-result
			if !ok {
				break
			}
			c.notifyDiagnosticResult(ctx, dia, s)
		}
		c.diagnosticFinish(ctx, dia)
	}
}

func (c *Coordinator) diagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) {
	for _, e := range c.exporters {
		c.logIfError(e.DiagnosticBegin(ctx, dia), "%s export diagnose begin", e.Meta().Name)
	}
}

func (c *Coordinator) diagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) {
	for _, e := range c.exporters {
		c.logIfError(e.DiagnosticFinish(ctx, dia), "%s export diagnose finish", e.Meta().Name)
	}
}

func (c *Coordinator) notifyDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) {
	for _, e := range c.exporters {
		c.logIfError(e.DiagnosticResult(ctx, dia, result), "%s export diagnose result", e.Meta().Name)
	}

}

func (c *Coordinator) logIfError(err error, format string, args ...interface{}) {
	if err != nil {
		c.logger.Errorf("%s failed : %v", fmt.Sprintf(format, args...), err.Error())
	}
}
