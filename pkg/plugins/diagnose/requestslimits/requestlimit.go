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
	d.Score = d.TotalScore
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
			d.diagnosePod(pod, d.Score/float64(len(pods.Items)))
		}
	}()
	return d.result
}

func (d Diagnostic) diagnosePod(pod v12.Pod, score float64) {
	for _, c := range pod.Spec.Containers {
		if c.Resources.Limits.Memory().IsZero() ||
			c.Resources.Limits.Cpu().IsZero() ||
			c.Resources.Requests.Memory().IsZero() ||
			c.Resources.Requests.Cpu().IsZero() {
			d.result <- &diagnose.Result{
				Level:    diagnose.HealthyLevelWarn,
				Name:     "Pods Requests Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Namespace, pod.Name),
				Desc:     d.Translator.Message("desc", nil),
				Score:    score,
				Proposal: d.Translator.Message("proposal", nil),
			}
			d.Score -= score
			return
		}
	}
}
