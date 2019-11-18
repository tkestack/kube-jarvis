package evaluate

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

// Result is the result of evaluation
type Result struct {
	Name string
	Desc string
}

// Evaluator knows how to evaluate all diagnostic results become one evaluation result
type Evaluator interface {
	// Param return core attributes
	Param() CreateParam
	// EvaDiagnosticResult evaluate one diagnostic result
	EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error
	// Result return a final evaluation result
	Result() *Result
}

// CreateParam contains core attributes of a Evaluator
type CreateParam struct {
	Logger logger.Logger
	Name   string
}

// Creator is a factory to create a Evaluator
type Creator func(c *CreateParam) Evaluator

// Creators store all registered Evaluator Creator
var Creators = map[string]Creator{}

// Add register a Evaluator Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
