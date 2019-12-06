package export

import (
	"context"

	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

// MetaData contains core attributes of a Exporter
type MetaData struct {
	plugins.CommonMetaData
}

// Meta return core MetaData
// this function can be use for struct implement Exporter interface
func (m *MetaData) Meta() MetaData {
	return *m
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
