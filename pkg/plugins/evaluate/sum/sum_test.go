package sum

import (
	"context"
	"fmt"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

func TestNewSumEva(t *testing.T) {
	s := NewEvaluator(&evaluate.CreateParam{
		Translator: translate.NewFake(),
	}).(*Evaluator)
	ctx := context.Background()
	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Level: diagnose.HealthyLevelRisk,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 2,
		Level: diagnose.HealthyLevelWarn,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 3,
		Level: diagnose.HealthyLevelSerious,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Error: fmt.Errorf("test"),
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnostic(ctx, example.NewDiagnostic(&diagnose.CreateParam{
		TotalScore: 6,
	})); err != nil {
		t.Fatalf(err.Error())
	}

	if s.TotalScore != 6 {
		t.Fatalf("total score should be 3")
	}

	if s.RiskScore != 1 {
		t.Fatalf("risk score should be 1")
	}

	if s.WarnScore != 2 {
		t.Fatalf("warn score should be 2")
	}

	if s.SeriousScore != 3 {
		t.Fatalf("warn score should be 3")
	}

	if s.ErrorTotal != 1 {
		t.Fatalf("error score should be 1")
	}

}
