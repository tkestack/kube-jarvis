package evaluate

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

type Result struct {
	Name string
	Desc string
}

type Evaluator interface {
	Param() CreateParam
	EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error
	Result() *Result
}

type CreateParam struct {
	Name string
}

type Creator func(c *CreateParam) Evaluator

var Creators = map[string]Creator{}

func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
