package evaluate

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

// MetaData contains core attributes of a Evaluator
type MetaData struct {
	plugins.CommonMetaData
}

// Meta return core MetaData
// this function can be use for struct implement Evaluator interface
func (m *MetaData) Meta() MetaData {
	return *m
}

// Result is the result of evaluation
type Result struct {
	// Title is the short description of evaluation result
	Name translate.Message
	// Desc is the full description of evaluation result
	Desc translate.Message
}

// Evaluator knows how to evaluate all diagnostic results become one evaluation result
type Evaluator interface {
	// Meta return core attributes
	Meta() MetaData
	// EvaDiagnosticResult evaluate one diagnostic result
	EvaDiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) error
	// EvaDiagnostic evaluate one diagnostic finish
	EvaDiagnostic(ctx context.Context, dia diagnose.Diagnostic) error
	// Result return a final evaluation result
	Result() *Result
}

// Factory create a new Evaluator
type Factory struct {
	// Creator is a factory function to create Evaluator
	Creator func(d *MetaData) Evaluator
	// SupportedClouds indicate what cloud providers will be supported of this evaluator
	SupportedClouds []string
}

// Factories store all registered Evaluator Creator
var Factories = map[string]Factory{}

// Add register a Evaluator Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}
