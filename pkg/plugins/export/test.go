package export

import (
	"context"
	"testing"

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

	d := example.NewDiagnostic(&diagnose.CreateParam{
		Translator: translate.NewFake(),
		Type:       "example",
		TotalScore: 10,
	})
	s := sum.NewEvaluator(&evaluate.CreateParam{
		Translator: translate.NewFake(),
		Type:       "sum",
		Name:       "sum",
	})

	// Diagnostic
	if err := exporter.DiagnosticBegin(ctx, d); err != nil {
		t.Fatalf(err.Error())
	}

	result := d.StartDiagnose(ctx)
	for {
		st, ok := <-result
		if !ok {
			break
		}

		if err := exporter.DiagnosticResult(ctx, st); err != nil {
			t.Fatalf(err.Error())
		}

		if err := s.EvaDiagnosticResult(ctx, st); err != nil {
			t.Fatalf(err.Error())
		}
	}

	if err := exporter.DiagnosticFinish(ctx, d); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnostic(ctx, d); err != nil {
		t.Fatalf(err.Error())
	}

	// Evaluation
	if err := exporter.EvaluationBegin(ctx, s); err != nil {
		t.Fatalf(err.Error())
	}

	if err := exporter.EvaluationResult(ctx, s.Result()); err != nil {
		t.Fatalf(err.Error())
	}

	if err := exporter.EvaluationFinish(ctx, s); err != nil {
		t.Fatalf(err.Error())
	}

	if err := exporter.CoordinateFinish(ctx); err != nil {
		t.Fatalf(err.Error())
	}
}
