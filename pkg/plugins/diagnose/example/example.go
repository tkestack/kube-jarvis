package example

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.CreateParam
	result  chan *diagnose.Result
	Message string `yaml:"message"`
}

// NewDiagnostic return a example diagnostic
func NewDiagnostic(config *diagnose.CreateParam) diagnose.Diagnostic {
	return &Diagnostic{
		result:      make(chan *diagnose.Result, 1000),
		CreateParam: config,
	}
}

// Param return core attributes
func (d *Diagnostic) Param() diagnose.CreateParam {
	return *d.CreateParam
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context) chan *diagnose.Result {
	go func() {
		defer close(d.result)
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelRisk,
			Name:     "sample",
			ObjName:  "sample-obj",
			Desc:     d.Message,
			Score:    10,
			Weight:   100,
			Proposal: "sample proposal",
		}
	}()
	return d.result
}
