package basic

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/coordinate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/Jarvis/pkg/logger"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
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
	for _, exp := range c.exporters {
		if err := exp.CoordinateBegin(ctx); err != nil {
			c.logger.Errorf("%s export coordinate begin failed : %v", err.Error())
		}
	}
}

func (c *Coordinator) finish(ctx context.Context) {
	for _, exp := range c.exporters {
		if err := exp.CoordinateFinish(ctx); err != nil {
			c.logger.Errorf("%s export coordinate finish failed : %v", err.Error())
		}
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
		if err := e.DiagnosticBegin(ctx, dia); err != nil {
			c.logger.Errorf("%s export diagnose begin failed : %v", err.Error())
		}
	}
}

func (c *Coordinator) diagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) {
	for _, e := range c.exporters {
		if err := e.DiagnosticFinish(ctx, dia); err != nil {
			c.logger.Errorf("%s export diagnose begin failed : %v", err.Error())
		}
	}
}

func (c *Coordinator) notifyDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) {
	for _, e := range c.exporters {
		if err := e.DiagnosticResult(ctx, result); err != nil {
			c.logger.Errorf("%s export diagnose result failed : %v", err.Error())
		}
	}

	for _, e := range c.evaluators {
		if err := e.EvaDiagnosticResult(ctx, result); err != nil {
			c.logger.Errorf("%s evaluator evaluate diagnose result failed : %v", err.Error())
		}
	}
}

func (c *Coordinator) evaluation(ctx context.Context) {
	for _, eva := range c.evaluators {
		result := eva.Result()
		for _, exp := range c.exporters {
			if err := exp.EvaluationBegin(ctx, eva); err != nil {
				c.logger.Errorf("%s export evaluation begin failed : %v", err.Error())
			}

			if err := exp.EvaluationResult(ctx, result); err != nil {
				c.logger.Errorf("%s export evaluation finish failed : %v", err.Error())
			}

			if err := exp.EvaluationFinish(ctx, eva); err != nil {
				c.logger.Errorf("%s export evaluation finish failed : %v", err.Error())
			}
		}
	}
}
