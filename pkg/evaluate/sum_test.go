package evaluate

import (
	"context"
	"fmt"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
)

func TestNewSumEva(t *testing.T) {
	s := NewSumEva()
	ctx := context.Background()
	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Level: diagnose.HealthyLevelRisk,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Level: diagnose.HealthyLevelPass,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Level: diagnose.HealthyLevelWarn,
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.EvaDiagnosticResult(ctx, &diagnose.Result{
		Score: 1,
		Error: fmt.Errorf("test"),
	}); err != nil {
		t.Fatalf(err.Error())
	}

	if s.TotalScore != 3 {
		t.Fatalf("total score should be 5")
	}

	if s.RiskScore != 1 {
		t.Fatalf("risk score should be 1")
	}

	if s.PassScore != 1 {
		t.Fatalf("pass score should be 1")
	}
	if s.WarnScore != 1 {
		t.Fatalf("warn score should be 1")
	}

	if s.ErrorTotal != 1 {
		t.Fatalf("error score should be 1")
	}

}
