package sys

import (
	"context"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "node-sys"
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
		for node, m := range d.param.Resources.Machines {
			curVal := m.SysCtl["net.ipv4.tcp_tw_reuse"]
			targetVal := "1"
			level := diagnose.HealthyLevelGood
			if curVal != targetVal {
				level = diagnose.HealthyLevelWarn
			}

			d.result <- &diagnose.Result{
				Level:   level,
				Title:   d.Translator.Message("kernel-para-title", nil),
				ObjName: node,
				Desc: d.Translator.Message("kernel-para-desc", map[string]interface{}{
					"Node":      node,
					"Name":      "net.ipv4.tcp_tw_reuse",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),

				Proposal: d.Translator.Message("kernel-para-proposal", map[string]interface{}{
					"Node":      node,
					"Name":      "net.ipv4.tcp_tw_reuse",
					"TargetVal": targetVal,
					"CurVal":    curVal,
				}),
			}
		}
	}()
	return d.result
}
