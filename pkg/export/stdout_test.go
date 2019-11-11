package export

import (
	"context"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
)

func TestNewStdout(t *testing.T) {
	s := NewStdout()
	ctx := context.Background()
	_ = s.CoordinateBegin(ctx)

	d := diagnose.NewSampleDiagnostic()
	sum := evaluate.NewSumEva()

	// Diagnostic
	_ = s.DiagnosticBegin(ctx, d)
	result := d.StartDiagnose(ctx, nil)
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
