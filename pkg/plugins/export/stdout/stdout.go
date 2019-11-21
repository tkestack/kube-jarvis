package stdout

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
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

// Param return core attributes
func (e *Exporter) Param() export.CreateParam {
	return *e.CreateParam
}

// CoordinateBegin export information about coordinator Run begin
func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	e.Logger.Infof("--------------- kube-jarvis [1.0] -----------------")
	return nil
}

// DiagnosticBegin export information about a Diagnostic begin
func (e *Exporter) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	e.Logger.Infof("[Diagnostic \"%s\" begin]", dia.Param().Name)
	return nil
}

// DiagnosticFinish export information about a Diagnostic finished
func (e *Exporter) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	e.Logger.Infof("[Diagnostic \"%s\" finish]", dia.Param().Name)
	return nil
}

// DiagnosticResult export information about one diagnose.Result
func (e *Exporter) DiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		e.Logger.Errorf(result.Error.Error())
	} else {
		e.Logger.Infof("[%s] [%s] [%s] [%s] [%s]", result.Level, result.Name, result.ObjName, result.Desc, result.Proposal)
	}
	return nil
}

// EvaluationBegin export information about a Evaluator begin
func (e *Exporter) EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error {
	e.Logger.Infof("[Evaluate \"%s\" begin]", eva.Param().Name)
	return nil
}

// EvaluationFinish export information about a Evaluator finish
func (e *Exporter) EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error {
	e.Logger.Infof("[Evaluate \"%s\" finish]", eva.Param().Name)
	return nil
}

// EvaluationResult export information about a Evaluator result
func (e *Exporter) EvaluationResult(ctx context.Context, result *evaluate.Result) error {
	e.Logger.Infof("[%s] [%s]", result.Name, result.Desc)
	return nil
}

// CoordinateFinish export information about coordinator Run finished
func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	e.Logger.Infof("--------------- Finish -----------------")
	return nil
}
