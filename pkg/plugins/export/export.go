package export

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
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

type Creator func() Exporter

var Creators = map[string]Creator{}

func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
