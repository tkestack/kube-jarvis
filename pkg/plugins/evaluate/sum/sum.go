package sum

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

// Evaluator sum all diagnostic result score with different healthy level
type Evaluator struct {
	*evaluate.CreateParam
	TotalScore   int
	WarnScore    int
	SeriousScore int
	RiskScore    int
	ErrorTotal   int
}

// NewEvaluator return a sum Evaluator
func NewEvaluator(param *evaluate.CreateParam) evaluate.Evaluator {
	return &Evaluator{
		CreateParam: param,
	}
}

// Param return core attributes
func (e *Evaluator) Param() evaluate.CreateParam {
	return *e.CreateParam
}

// EvaDiagnosticResult evaluate one diagnostic result
func (e *Evaluator) EvaDiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		e.ErrorTotal++
	} else {
		switch result.Level {
		case diagnose.HealthyLevelWarn:
			e.WarnScore += result.Score
		case diagnose.HealthyLevelRisk:
			e.RiskScore += result.Score
		case diagnose.HealthyLevelSerious:
			e.SeriousScore += result.Score
		}
		e.TotalScore += result.Score
	}
	return nil
}

// Result return a final evaluation result
func (e *Evaluator) Result() *evaluate.Result {
	return &evaluate.Result{
		Name: "score statistics",
		Desc: e.Translator.Message("result", map[string]interface{}{
			"TotalScore":   e.TotalScore,
			"WarnScore":    e.WarnScore,
			"RiskScore":    e.RiskScore,
			"SeriousScore": e.SeriousScore,
			"ErrorTotal":   e.ErrorTotal,
		}),
	}
}
