package example

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result  chan *diagnose.Result
	Message string `yaml:"message"`
}

// NewDiagnostic return a example diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:   make(chan *diagnose.Result, 1000),
		MetaData: meta,
	}
}

// Meta return core attributes
func (d *Diagnostic) Meta() diagnose.MetaData {
	return *d.MetaData
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context) chan *diagnose.Result {
	go func() {
		defer close(d.result)
		d.result <- &diagnose.Result{
			Level:   diagnose.HealthyLevelRisk,
			Name:    "example",
			ObjName: "example-obj",
			Desc: d.Translator.Message("message", map[string]interface{}{
				"Mes": d.Message,
			}),
			Score:    d.TotalScore,
			Proposal: d.Translator.Message("proposal", nil),
		}
	}()
	return d.result
}
