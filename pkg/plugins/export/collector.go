package export

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"
	"gopkg.in/yaml.v2"
)

// DiagnosticResultItem collect one diagnostic and it's results
type DiagnosticResultItem struct {
	Type       string
	Name       string
	Score      float64
	TotalScore float64
	Results    []diagnose.Result
}

// EvaluationResultItem collect one evaluator and it's result
type EvaluationResultItem struct {
	Type   string
	Name   string
	Result evaluate.Result
}

// Collector just collect diagnostic results and evaluation results
type Collector struct {
	Diagnostics []*DiagnosticResultItem
	Evaluations []*EvaluationResultItem
}

// Param return core attributes
func (c *Collector) Param() CreateParam {
	return CreateParam{}
}

// CoordinateBegin export information about coordinator Run begin
func (c *Collector) CoordinateBegin(ctx context.Context) error {
	return nil
}

// DiagnosticBegin export information about a Diagnostic begin
func (c *Collector) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	param := dia.Param()
	c.Diagnostics = append(c.Diagnostics, &DiagnosticResultItem{
		Type:       param.Type,
		Name:       param.Name,
		TotalScore: param.TotalScore,
	})
	return nil
}

// DiagnosticResult export information about one diagnose.Result
func (c *Collector) DiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	dLen := len(c.Diagnostics)
	if dLen == 0 {
		return fmt.Errorf("no diagnostic found")
	}

	c.Diagnostics[dLen-1].Results = append(c.Diagnostics[dLen-1].Results, *result)
	return nil
}

// DiagnosticFinish export information about a Diagnostic finished
func (c *Collector) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	dLen := len(c.Diagnostics)
	if dLen == 0 {
		return fmt.Errorf("no diagnostic found")
	}
	c.Diagnostics[dLen-1].Score = dia.Param().Score
	return nil
}

// EvaluationBegin export information about a Evaluator begin
func (c *Collector) EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error {
	param := eva.Param()
	c.Evaluations = append(c.Evaluations, &EvaluationResultItem{
		Type: param.Type,
		Name: param.Name,
	})
	return nil
}

// EvaluationResult export information about a Evaluator result
func (c *Collector) EvaluationResult(ctx context.Context, result *evaluate.Result) error {
	eLen := len(c.Evaluations)
	if eLen == 0 {
		return fmt.Errorf("no evaluations found")
	}

	c.Evaluations[eLen-1].Result = *result
	return nil
}

// EvaluationFinish export information about a Evaluator finish
func (c *Collector) EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error {
	return nil
}

// CoordinateFinish export information about coordinator Run finished
func (c *Collector) CoordinateFinish(ctx context.Context) error {
	return nil
}

// Marshal marshal Collected results to byte data according to format
// format can be : "json" , "yaml"
func (c *Collector) Marshal(format string) ([]byte, error) {
	switch format {
	case "json":
		return json.Marshal(c)
	case "yaml":
		return yaml.Marshal(c)
	}

	return nil, fmt.Errorf("unknow format")
}

// Unmarshal unmarshal data to Collector
// format can be : "json" , "yaml"
func (c *Collector) Unmarshal(format string, data []byte) error {
	switch format {
	case "json":
		return json.Unmarshal(data, c)
	case "yaml":
		return yaml.Unmarshal(data, c)
	}
	return fmt.Errorf("unknow format")
}
