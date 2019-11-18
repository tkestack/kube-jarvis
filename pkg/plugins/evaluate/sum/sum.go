package sum

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

type Evaluator struct {
	*evaluate.CreateParam
	TotalScore int
	PassScore  int
	WarnScore  int
	RiskScore  int
	ErrorTotal int
}

func NewEvaluator(param *evaluate.CreateParam) evaluate.Evaluator {
	return &Evaluator{
		CreateParam: param,
	}
}

func (e *Evaluator) Param() evaluate.CreateParam {
	return *e.CreateParam
}

func (e *Evaluator) EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		e.ErrorTotal++
	} else {
		switch result.Level {
		case diagnose.HealthyLevelPass:
			e.PassScore += result.Score
		case diagnose.HealthyLevelWarn:
			e.WarnScore += result.Score
		case diagnose.HealthyLevelRisk:
			e.RiskScore += result.Score
		}
		e.TotalScore += result.Score
	}
	return nil
}

func (e *Evaluator) Result() *evaluate.Result {
	return &evaluate.Result{
		Name: "score statistics",
		Desc: fmt.Sprintf("TotalScore : %d, PassScore : %d, WarnScore : %d, RiskScore : %d, ErrorTotal : %d",
			e.TotalScore, e.PassScore, e.WarnScore, e.RiskScore, e.ErrorTotal),
	}
}
