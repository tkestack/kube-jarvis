package stdout

import (
	"context"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	sum2 "github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate/sum"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"
)

func TestNewStdout(t *testing.T) {
	s := NewExporter(&export.CreateParam{}).(*Exporter)
	ctx := context.Background()
	_ = s.CoordinateBegin(ctx)

	d := example.NewDiagnostic(&diagnose.CreateParam{
		Translator: translate.NewFake(),
		Score:      10,
		Weight:     10,
	})
	sum := sum2.NewEvaluator(&evaluate.CreateParam{
		Translator: translate.NewFake(),
		Name:       "sum",
	})

	// Diagnostic
	_ = s.DiagnosticBegin(ctx, d)
	result := d.StartDiagnose(ctx)
	for {
		st, ok := <-result
		if !ok {
			break
		}
		_ = s.DiagnosticResult(ctx, st)
		_ = sum.EvaDiagnosticResult(ctx, st)
	}
	_ = s.DiagnosticFinish(ctx, d)

	// Evaluation
	_ = s.EvaluationBegin(ctx, sum)
	_ = s.EvaluationResult(ctx, sum.Result())
	_ = s.EvaluationFinish(ctx, sum)

	_ = s.CoordinateFinish(ctx)
}
