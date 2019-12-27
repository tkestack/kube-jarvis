package coordinate

import (
	"context"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

type FakeCoordinator struct {
	RunFunc func(ctx context.Context)
}

// Complete check and complete config items
func (f *FakeCoordinator) Complete() error {
	return nil
}

// AddDiagnostic add a diagnostic to Coordinator
func (f *FakeCoordinator) AddDiagnostic(dia diagnose.Diagnostic) {

}

// AddExporter add a Exporter to Coordinator
func (f *FakeCoordinator) AddExporter(exporter export.Exporter) {

}

// Run will do all diagnostics, evaluations, then export it by exporters
func (f *FakeCoordinator) Run(ctx context.Context) {
	if f.RunFunc != nil {
		f.RunFunc(ctx)
	}
}

func (f *FakeCoordinator) Progress() *plugins.Progress {
	return nil
}
