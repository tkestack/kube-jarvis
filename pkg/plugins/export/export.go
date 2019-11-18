package export

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

// Exporter export all steps and results with special way or special format
type Exporter interface {
	Param() CreateParam
	CoordinateBegin(ctx context.Context) error

	DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error
	DiagnosticResult(ctx context.Context, result *diagnose.Result) error
	DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error

	EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error
	EvaluationResult(ctx context.Context, result *evaluate.Result) error
	EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error

	CoordinateFinish(ctx context.Context) error
}

// CreateParam contains core attributes of a Exporter
type CreateParam struct {
	Logger logger.Logger
	Name   string
}

// Creator is a factory to create a Exporter
type Creator func(c *CreateParam) Exporter

// Creators store all registered Exporter Creator
var Creators = map[string]Creator{}

// Add register a Exporter Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
