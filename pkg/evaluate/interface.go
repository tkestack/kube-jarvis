package evaluate

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
)

type Result struct {
	Name string
	Desc string
}

type Evaluator interface {
	Name() string
	EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error
	Result() *Result
}
