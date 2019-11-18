package stdout

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/logger"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

// Exporter just print information to logger with a simple format
type Exporter struct {
	*export.CreateParam
	Logger logger.Logger
}

// NewExporter return a stdout Exporter
func NewExporter(p *export.CreateParam) export.Exporter {
	return &Exporter{
		Logger:      logger.NewLogger(),
		CreateParam: p,
	}
}

func (e *Exporter) Param() export.CreateParam {
	return *e.CreateParam
}

func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	e.Logger.Infof("--------------- Jarvis [1.0] -----------------")
	return nil
}

func (e *Exporter) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	e.Logger.Infof("[Diagnostic \"%s\" begin]", dia.Param().Name)
	return nil
}

func (e *Exporter) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	e.Logger.Infof("[Diagnostic \"%s\" finish]", dia.Param().Name)
	return nil
}

func (e *Exporter) DiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		e.Logger.Errorf(result.Error.Error())
	} else {
		e.Logger.Infof("[%s] [%s] [%s] [%s] [%s]", result.Level, result.Name, result.ObjName, result.Desc, result.Proposal)
	}
	return nil
}

func (e *Exporter) EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error {
	e.Logger.Infof("[Evaluate \"%s\" begin]", eva.Param().Name)
	return nil
}

func (e *Exporter) EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error {
	e.Logger.Infof("[Evaluate \"%s\" finish]", eva.Param().Name)
	return nil
}

func (e *Exporter) EvaluationResult(ctx context.Context, result *evaluate.Result) error {
	e.Logger.Infof("[%s] [%s]", result.Name, result.Desc)
	return nil
}

func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	e.Logger.Infof("--------------- Finish -----------------")
	return nil
}
