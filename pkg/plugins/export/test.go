package export

import (
	"context"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate/sum"
	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
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
		TotalScore: 10,
	})
	s := sum.NewEvaluator(&evaluate.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
			Type:       "sum",
			Name:       "sum",
		},
	})

	// Diagnostic
	fatalIfError(t, exporter.DiagnosticBegin(ctx, d))
	result := d.StartDiagnose(ctx)
	for {
		st, ok := <-result
		if !ok {
			break
		}

		fatalIfError(t, exporter.DiagnosticResult(ctx, d, st))
		fatalIfError(t, s.EvaDiagnosticResult(ctx, d, st))
	}

	fatalIfError(t, exporter.DiagnosticFinish(ctx, d))
	fatalIfError(t, s.EvaDiagnostic(ctx, d))
	// Evaluation
	fatalIfError(t, exporter.EvaluationBegin(ctx, s))
	fatalIfError(t, exporter.EvaluationResult(ctx, s, s.Result()))
	fatalIfError(t, exporter.EvaluationFinish(ctx, s))
	fatalIfError(t, exporter.CoordinateFinish(ctx))
}

func fatalIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}
