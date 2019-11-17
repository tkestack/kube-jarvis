package stdout

import (
	"context"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"

	sum2 "github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate/sum"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/example"
)

func TestNewStdout(t *testing.T) {
	s := NewExporter(&export.CreateParam{}).(*Exporter)
	ctx := context.Background()
	_ = s.CoordinateBegin(ctx)

	d := example.NewDiagnostic(&diagnose.CreateParam{
		Score:  10,
		Weight: 10,
	})
	sum := sum2.NewEvaluator(&evaluate.CreateParam{
		Name: "sum",
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
