package evaluate

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
)

type SumEva struct {
	NameDesc   string
	TotalScore int
	PassScore  int
	WarnScore  int
	RiskScore  int
	ErrorTotal int
}

func NewSumEva() *SumEva {
	return &SumEva{
		NameDesc: "SumEva",
	}
}

func (s *SumEva) Name() string {
	return s.NameDesc
}

func (s *SumEva) EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		s.ErrorTotal++
	} else {
		switch result.Level {
		case diagnose.HealthyLevelPass:
			s.PassScore += result.Score
		case diagnose.HealthyLevelWarn:
			s.WarnScore += result.Score
		case diagnose.HealthyLevelRisk:
			s.RiskScore += result.Score
		}
		s.TotalScore += result.Score
	}
	return nil
}

func (s *SumEva) Result() *Result {
	return &Result{
		Name: "score statistics",
		Desc: fmt.Sprintf("TotalScore : %d, PassScore : %d, WarnScore : %d, RiskScore : %d, ErrorTotal : %d",
			s.TotalScore, s.PassScore, s.WarnScore, s.RiskScore, s.ErrorTotal),
	}
}
