package export

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

// MetaData contains core attributes of a Exporter
type MetaData struct {
	plugins.CommonMetaData
}

// Exporter export all steps and results with special way or special format
type Exporter interface {
	// Meta return core attributes
	Meta() MetaData
	// CoordinateBegin export information about coordinator Run begin
	CoordinateBegin(ctx context.Context) error
	// DiagnosticBegin export information about a Diagnostic begin
	DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error
	// DiagnosticResult export information about one diagnose.Result
	DiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) error
	// DiagnosticFinish export information about a Diagnostic finished
	DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error
	// EvaluationBegin export information about a Evaluator begin
	EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error
	// EvaluationResult export information about a Evaluator result
	EvaluationResult(ctx context.Context, eva evaluate.Evaluator, result *evaluate.Result) error
	// EvaluationFinish export information about a Evaluator finish
	EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error
	// CoordinateFinish export information about coordinator Run finished
	CoordinateFinish(ctx context.Context) error
}

// Factory create a new Exporter
type Factory struct {
	// Creator is a factory function to create Exporter
	Creator func(d *MetaData) Exporter
	// SupportedClouds indicate what cloud providers will be supported of this exporter
	SupportedClouds []string
}

// Factories store all registered Exporter Creator
var Factories = map[string]Factory{}

// Add register a Exporter Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}

// IsSupported return true if cloud type is supported by Exporter
func (f *Factory) IsSupported(cloud string) bool {
	if len(f.SupportedClouds) == 0 {
		return true
	}

	for _, c := range f.SupportedClouds {
		if c == cloud {
			return true
		}
	}
	return false
}
