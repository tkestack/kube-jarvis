package basic

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

// Coordinator Coordinate diagnostics,exporters,evaluators with simple way
type Coordinator struct {
	logger      logger.Logger
	diagnostics []diagnose.Diagnostic
	exporters   []export.Exporter
	evaluators  []evaluate.Evaluator
}

// NewCoordinator return a default Coordinator
func NewCoordinator(logger logger.Logger) coordinate.Coordinator {
	return &Coordinator{
		logger: logger,
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

// AddEvaluate add a evaluate to Coordinator
func (c *Coordinator) AddEvaluate(evaluate evaluate.Evaluator) {
	c.evaluators = append(c.evaluators, evaluate)
}

// Run will do all diagnostics, evaluations, then export it by exporters
func (c *Coordinator) Run(ctx context.Context) {
	c.begin(ctx)
	c.diagnostic(ctx)
	c.evaluation(ctx)
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
		result := dia.StartDiagnose(ctx)
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
	for _, e := range c.evaluators {
		c.logIfError(e.EvaDiagnostic(ctx, dia), "%s evaluate diagnose finish", e.Meta().Name)
	}
}

func (c *Coordinator) notifyDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) {
	for _, e := range c.exporters {
		c.logIfError(e.DiagnosticResult(ctx, dia, result), "%s export diagnose result", e.Meta().Name)
	}

	for _, e := range c.evaluators {
		c.logIfError(e.EvaDiagnosticResult(ctx, dia, result), "%s evaluator evaluate diagnose result begin", e.Meta().Name)
	}
}

func (c *Coordinator) evaluation(ctx context.Context) {
	for _, eva := range c.evaluators {
		result := eva.Result()
		for _, exp := range c.exporters {
			expName := exp.Meta().Name
			c.logIfError(exp.EvaluationBegin(ctx, eva), "%s export evaluation begin", expName)
			c.logIfError(exp.EvaluationResult(ctx, eva, result), "%s export evaluation result", expName)
			c.logIfError(exp.EvaluationFinish(ctx, eva), "%s export evaluation finish", expName)
		}
	}
}

func (c *Coordinator) logIfError(err error, format string, args ...interface{}) {
	if err != nil {
		c.logger.Errorf("%s failed : %v", fmt.Sprintf(format, args...), err.Error())
	}
}
