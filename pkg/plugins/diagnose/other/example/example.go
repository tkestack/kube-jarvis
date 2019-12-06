package example

import (
	"context"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "example"
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

// Init do initialization
func (d *Diagnostic) Init() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) chan *diagnose.Result {
	go func() {
		defer close(d.result)
		d.result <- &diagnose.Result{
			Level:   diagnose.HealthyLevelRisk,
			Title:   "example",
			ObjName: "example-obj",
			Desc: d.Translator.Message("message", map[string]interface{}{
				"Mes": d.Message,
			}),
			Proposal: d.Translator.Message("proposal", nil),
		}
	}()
	return d.result
}
