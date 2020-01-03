package status

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubectl/pkg/describe"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "node-status"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a node status diagnostic
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

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer close(d.result)
		for _, node := range d.param.Resources.Nodes.Items {
			isMaster := false
			isHealth := true
			levelGood := diagnose.HealthyLevelGood
			levelBad := diagnose.HealthyLevelRisk
			if node.Labels[describe.LabelNodeRolePrefix+"master"] == "true" {
				isMaster = true
				levelBad = diagnose.HealthyLevelSerious
			}
			for _, cond := range node.Status.Conditions {
				if cond.Status == v1.ConditionUnknown {
					d.uploadResult(isMaster, node.Name, v1.NodeReady, cond.Status, levelBad)
					isHealth = false
					break
				}
				if (cond.Type == v1.NodeReady && cond.Status != v1.ConditionTrue) || (cond.Type != v1.NodeReady && cond.Status == v1.ConditionTrue) {
					d.uploadResult(isMaster, node.Name, cond.Type, cond.Status, levelBad)
					isHealth = false
					break
				}
			}
			if isHealth {
				d.uploadResult(isMaster, node.Name, v1.NodeReady, v1.ConditionTrue, levelGood)
			}
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) uploadResult(isMaster bool, name string, typ v1.NodeConditionType, status v1.ConditionStatus, level diagnose.HealthyLevel) {
	resource := typ
	prefix := "node"
	goodFlag := ""
	if isMaster {
		prefix = "master"
	}
	if level == diagnose.HealthyLevelGood {
		goodFlag = "good-"
	}

	title := d.Translator.Message(prefix+"-status-title", nil)
	desc := d.Translator.Message(prefix+"-status-"+goodFlag+"desc", map[string]interface{}{
		"Node":   name,
		"Type":   typ,
		"Status": status,
	})
	proposal := d.Translator.Message(prefix+"-status-"+goodFlag+"proposal", map[string]interface{}{
		"Node":     name,
		"Resource": resource,
	})

	d.result <- &diagnose.Result{
		Level:    level,
		Title:    title,
		ObjName:  name,
		Desc:     desc,
		Proposal: proposal,
	}
}
