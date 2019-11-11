package export

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
)

type Exporter interface {
	CoordinateBegin(ctx context.Context) error

	DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error
	DiagnosticResult(ctx context.Context, result *diagnose.Result) error
	DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error

	EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error
	EvaluationResult(ctx context.Context, result *evaluate.Result) error
	EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error

	CoordinateFinish(ctx context.Context) error
}
