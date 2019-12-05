package apiserver

import (
	"context"
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "master-apiserver"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
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
	d.param = &param
	go func() {
		defer close(d.result)
		for _, m := range d.param.Resources.CoreComponents[cluster.ComponentApiserver] {
			curVal := m.Args["allow-privileged"]
			targetVal := "false"
			level := diagnose.HealthyLevelGood
			if curVal != targetVal {
				level = diagnose.HealthyLevelWarn
			}

			d.result <- &diagnose.Result{
				Level:   level,
				Title:   d.Translator.Message("apiserver-para-title", nil),
				ObjName: fmt.Sprintf("%s(%s)", cluster.ComponentApiserver, m.Node),
				Desc: d.Translator.Message("apiserver-para-desc", map[string]interface{}{
					"Node":      m.Node,
					"Name":      "allow-privileged",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),

				Proposal: d.Translator.Message("apiserver-para-proposal", map[string]interface{}{
					"Node":      m.Node,
					"Name":      "allow-privileged",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),
			}
		}
	}()
	return d.result
}
