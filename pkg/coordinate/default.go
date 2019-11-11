package coordinate

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/RayHuangCN/Jarvis/pkg/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/export"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/logger"
)

type DefCoordination struct {
	logger      logger.Logger
	cli         kubernetes.Interface
	diagnostics []diagnose.Diagnostic
	exporters   []export.Exporter
	evaluators  []evaluate.Evaluator
}

func NewDefault(logger logger.Logger, cli kubernetes.Interface) *DefCoordination {
	return &DefCoordination{
		logger: logger,
		cli:    cli,
	}
}

func (d *DefCoordination) AddDiagnostic(dia diagnose.Diagnostic) {
	d.diagnostics = append(d.diagnostics, dia)
}

func (d *DefCoordination) AddExporter(exporter export.Exporter) {
	d.exporters = append(d.exporters, exporter)
}

func (d *DefCoordination) AddEvaluate(evaluate evaluate.Evaluator) {
	d.evaluators = append(d.evaluators, evaluate)
}

func (d *DefCoordination) Run(ctx context.Context) {
	d.begin(ctx)
	d.diagnostic(ctx)
	d.evaluation(ctx)
	d.finish(ctx)
}

func (d *DefCoordination) begin(ctx context.Context) {
	for _, exp := range d.exporters {
		if err := exp.CoordinateBegin(ctx); err != nil {
			d.logger.Errorf("%s export coordinate begin failed : %v", err.Error())
		}
	}
}

func (d *DefCoordination) finish(ctx context.Context) {
	for _, exp := range d.exporters {
		if err := exp.CoordinateFinish(ctx); err != nil {
			d.logger.Errorf("%s export coordinate finish failed : %v", err.Error())
		}
	}
}

func (d *DefCoordination) diagnostic(ctx context.Context) {
	for _, dia := range d.diagnostics {
		d.diagnosticBegin(ctx, dia)
		result := dia.StartDiagnose(ctx, d.cli)
		for {
			s, ok := <-result
			if !ok {
				break
			}
			d.notifyDiagnosticResult(ctx, dia, s)
		}
		d.diagnosticFinish(ctx, dia)
	}
}

func (d *DefCoordination) diagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) {
	for _, e := range d.exporters {
		if err := e.DiagnosticBegin(ctx, dia); err != nil {
			d.logger.Errorf("%s export diagnose begin failed : %v", err.Error())
		}
	}
}

func (d *DefCoordination) diagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) {
	for _, e := range d.exporters {
		if err := e.DiagnosticFinish(ctx, dia); err != nil {
			d.logger.Errorf("%s export diagnose begin failed : %v", err.Error())
		}
	}
}

func (d *DefCoordination) notifyDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) {
	for _, e := range d.exporters {
		if err := e.DiagnosticResult(ctx, result); err != nil {
			d.logger.Errorf("%s export diagnose result failed : %v", err.Error())
		}
	}

	for _, e := range d.evaluators {
		if err := e.EvaDiagnosticResult(ctx, result); err != nil {
			d.logger.Errorf("%s evaluator evaluate diagnose result failed : %v", err.Error())
		}
	}
}

func (d *DefCoordination) evaluation(ctx context.Context) {
	for _, eva := range d.evaluators {
		result := eva.Result()
		for _, exp := range d.exporters {
			if err := exp.EvaluationBegin(ctx, eva); err != nil {
				d.logger.Errorf("%s export evaluation begin failed : %v", err.Error())
			}

			if err := exp.EvaluationResult(ctx, result); err != nil {
				d.logger.Errorf("%s export evaluation finish failed : %v", err.Error())
			}

			if err := exp.EvaluationFinish(ctx, eva); err != nil {
				d.logger.Errorf("%s export evaluation finish failed : %v", err.Error())
			}
		}
	}
}
