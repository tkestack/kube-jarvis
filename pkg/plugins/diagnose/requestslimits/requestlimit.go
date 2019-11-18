package requestslimits

import (
	"context"
	"fmt"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"

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

func (d *Diagnostic) Param() diagnose.CreateParam {
	return *d.CreateParam
}

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
		if c.Resources.Limits.Memory().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Name:     "Pods Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no memory limits",
				Score:    1,
				Weight:   d.Weight,
				Proposal: "container should set memory limits",
			}
		}

		if c.Resources.Limits.Cpu().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Name:     "Pods Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no cpu limits",
				Score:    1,
				Weight:   d.Weight,
				Proposal: "any container should set cpu limits",
			}
		}

		if c.Resources.Requests.Memory().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Name:     "Pods Requests",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no memory requests",
				Score:    1,
				Weight:   d.Weight,
				Proposal: "container should set memory requests",
			}
		}

		if c.Resources.Requests.Cpu().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Name:     "Pods Requests",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no cpu requests",
				Score:    1,
				Weight:   d.Weight,
				Proposal: "any container should set cpu requests",
			}
		}
	}
}
