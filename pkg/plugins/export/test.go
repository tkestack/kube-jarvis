package export

import (
	"context"
	"testing"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/other/example"
	"tkestack.io/kube-jarvis/pkg/translate"
)

// RunExporterTest is a tool function for exporter testing
func RunExporterTest(t *testing.T, exporter Exporter) {
	ctx := context.Background()
	_ = exporter.CoordinateBegin(ctx)

	d := example.NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
			Type:       "example",
		},
	})
	// Diagnostic
	fatalIfError(t, exporter.DiagnosticBegin(ctx, d))
	result := d.StartDiagnose(ctx, diagnose.StartDiagnoseParam{})
	for {
		st, ok := <-result
		if !ok {
			break
		}

		fatalIfError(t, exporter.DiagnosticResult(ctx, d, st))
	}

	fatalIfError(t, exporter.DiagnosticFinish(ctx, d))
	// Evaluation
	fatalIfError(t, exporter.CoordinateFinish(ctx))
}

func fatalIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}
