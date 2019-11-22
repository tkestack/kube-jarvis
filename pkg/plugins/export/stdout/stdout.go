package stdout

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"
	"github.com/fatih/color"
)

// Exporter just print information to logger with a simple format
type Exporter struct {
	*export.CreateParam
}

// NewExporter return a stdout Exporter
func NewExporter(p *export.CreateParam) export.Exporter {
	return &Exporter{
		CreateParam: p,
	}
}

// Param return core attributes
func (e *Exporter) Param() export.CreateParam {
	return *e.CreateParam
}

// CoordinateBegin export information about coordinator Run begin
func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	fmt.Println("===================================================================")
	fmt.Println("                       kube-jarivs                                 ")
	fmt.Println("===================================================================")
	return nil
}

// DiagnosticBegin export information about a Diagnostic begin
func (e *Exporter) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	fmt.Println("Diagnostic report")
	fmt.Printf("    Type : %s\n", dia.Param().Type)
	fmt.Printf("    Name : %s\n", dia.Param().Name)
	fmt.Printf("- ----- results ----------------\n")
	return nil
}

// DiagnosticFinish export information about a Diagnostic finished
func (e *Exporter) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	fmt.Printf("Diagnostic Score : %.2f/%.2f\n", dia.Param().Score, dia.Param().TotalScore)
	fmt.Println("===================================================================")
	return nil
}

// DiagnosticResult export information about one diagnose.Result
func (e *Exporter) DiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		color.HiRed("[!!ERROR] %s\n", result.Error.Error())
	} else {
		var pt func(format string, a ...interface{})
		switch result.Level {
		case diagnose.HealthyLevelWarn:
			pt = color.Yellow
		case diagnose.HealthyLevelRisk:
			pt = color.Red
		case diagnose.HealthyLevelSerious:
			pt = color.HiRed
		default:
			pt = func(format string, a ...interface{}) {
				fmt.Printf(format, a...)
			}
		}
		pt("[%s] %s -> %s\n", result.Level, result.Name, result.ObjName)
		pt("    Score : -%.2f\n", result.Score)
		pt("    Describe : %s\n", result.Desc)
		pt("    Proposal : %s\n", result.Proposal)
	}
	fmt.Printf("- -----------------------------\n")
	return nil
}

// EvaluationBegin export information about a Evaluator begin
func (e *Exporter) EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error {
	fmt.Println("Evaluation report")
	fmt.Printf("    Type : %s\n", eva.Param().Type)
	fmt.Printf("    Name : %s\n", eva.Param().Name)
	fmt.Printf("- ----- result -----------------\n")
	return nil
}

// EvaluationFinish export information about a Evaluator finish
func (e *Exporter) EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error {
	fmt.Println("===================================================================")
	return nil
}

// EvaluationResult export information about a Evaluator result
func (e *Exporter) EvaluationResult(ctx context.Context, result *evaluate.Result) error {
	fmt.Printf("[%s]\n", result.Name)
	fmt.Printf("    Describe : %s\n", result.Desc)
	return nil
}

// CoordinateFinish export information about coordinator Run finished
func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	fmt.Println("===================================================================")
	return nil
}
