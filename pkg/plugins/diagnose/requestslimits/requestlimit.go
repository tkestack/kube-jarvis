package requestslimits

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Diagnostic report the healthy of pods's resources requests limits configuration
type Diagnostic struct {
	*diagnose.CreateParam
	result chan *diagnose.Result
}

// NewDiagnostic return a requests-limits Diagnostic
func NewDiagnostic(config *diagnose.CreateParam) diagnose.Diagnostic {
	return &Diagnostic{
		CreateParam: config,
		result:      make(chan *diagnose.Result, 1000),
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
		defer func() {
			if err := recover(); err != nil {
				d.result <- &diagnose.Result{
					Error: fmt.Errorf("%v", err),
				}
			}
		}()

		pods, err := d.Cli.CoreV1().Pods("").List(v1.ListOptions{})
		if err != nil {
			d.result <- &diagnose.Result{
				Error: err,
			}
			return
		}

		for _, pod := range pods.Items {
			d.diagnosePod(pod)
		}
	}()
	return d.result
}

func (d Diagnostic) diagnosePod(pod v12.Pod) {
	for _, c := range pod.Spec.Containers {
		health := true
		if c.Resources.Limits.Memory().IsZero() {
			health = false
			d.result <- &diagnose.Result{
				Level:   diagnose.HealthyLevelWarn,
				Name:    "Pods Limits",
				ObjName: fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc: d.Translator.Message("desc", map[string]interface{}{
					"Resource": "memory",
					"Type":     "limits",
				}),
				Score:  d.Score,
				Weight: d.Weight,
				Proposal: d.Translator.Message("proposal", map[string]interface{}{
					"Resource": "memory",
					"Type":     "limits",
				}),
			}
		}

		if c.Resources.Limits.Cpu().IsZero() {
			health = false
			d.result <- &diagnose.Result{
				Level:   diagnose.HealthyLevelWarn,
				Name:    "Pods Limits",
				ObjName: fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc: d.Translator.Message("desc", map[string]interface{}{
					"Resource": "cpu",
					"Type":     "limits",
				}),
				Score:  d.Score,
				Weight: d.Weight,
				Proposal: d.Translator.Message("proposal", map[string]interface{}{
					"Resource": "memory",
					"Type":     "limits",
				}),
			}
		}

		if c.Resources.Requests.Memory().IsZero() {
			health = false
			d.result <- &diagnose.Result{
				Level:   diagnose.HealthyLevelWarn,
				Name:    "Pods Requests",
				ObjName: fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc: d.Translator.Message("desc", map[string]interface{}{
					"Resource": "memory",
					"Type":     "requests",
				}),
				Score:  d.Score,
				Weight: d.Weight,
				Proposal: d.Translator.Message("proposal", map[string]interface{}{
					"Resource": "memory",
					"Type":     "requests",
				}),
			}
		}

		if c.Resources.Requests.Cpu().IsZero() {
			health = false
			d.result <- &diagnose.Result{
				Level:   diagnose.HealthyLevelWarn,
				Name:    "Pods Requests",
				ObjName: fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc: d.Translator.Message("desc", map[string]interface{}{
					"Resource": "cpu",
					"Type":     "requests",
				}),
				Score:  d.Score,
				Weight: d.Weight,
				Proposal: d.Translator.Message("proposal", map[string]interface{}{
					"Resource": "cpu",
					"Type":     "requests",
				}),
			}
		}

		if health {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelPass,
				Name:     "Pods resources",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     d.Translator.Message("passDesc", nil),
				Score:    d.Score,
				Weight:   d.Weight,
				Proposal: "",
			}
		}

	}
}
