package diagnose

import (
	"context"
	"fmt"

	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type RequestLimitDiagnostic struct {
	result   chan *Result
	Weight   int
	NameDesc string
}

func NewRequestLimitDiagnostic() *RequestLimitDiagnostic {
	return &RequestLimitDiagnostic{
		result:   make(chan *Result, 1000),
		Weight:   1,
		NameDesc: "RequestLimitDiagnostic",
	}
}

func (r *RequestLimitDiagnostic) Name() string {
	return r.NameDesc
}

func (r *RequestLimitDiagnostic) StartDiagnose(ctx context.Context, cli kubernetes.Interface) chan *Result {
	go func() {
		defer close(r.result)
		defer func() {
			if err := recover(); err != nil {
				r.result <- &Result{
					Error: fmt.Errorf("%v", err),
				}
			}
		}()

		pods, err := cli.CoreV1().Pods("").List(v1.ListOptions{})
		if err != nil {
			r.result <- &Result{
				Error: err,
			}
			return
		}

		for _, pod := range pods.Items {
			r.diagnosePod(pod)
		}
	}()
	return r.result
}

func (r RequestLimitDiagnostic) diagnosePod(pod v12.Pod) {
	for _, c := range pod.Spec.Containers {
		if c.Resources.Limits.Memory().IsZero() {
			r.result <- &Result{
				Level:    HealthyLevelWarn,
				Name:     "Pods Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no memory limits",
				Score:    1,
				Weight:   r.Weight,
				Proposal: "container should set memory limits",
			}
		}

		if c.Resources.Limits.Cpu().IsZero() {
			r.result <- &Result{
				Level:    HealthyLevelWarn,
				Name:     "Pods Limits",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no cpu limits",
				Score:    1,
				Weight:   r.Weight,
				Proposal: "any container should set cpu limits",
			}
		}

		if c.Resources.Requests.Memory().IsZero() {
			r.result <- &Result{
				Level:    HealthyLevelWarn,
				Name:     "Pods Requests",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no memory requests",
				Score:    1,
				Weight:   r.Weight,
				Proposal: "container should set memory requests",
			}
		}

		if c.Resources.Requests.Cpu().IsZero() {
			r.result <- &Result{
				Level:    HealthyLevelWarn,
				Name:     "Pods Requests",
				ObjName:  fmt.Sprintf("%s:%s", pod.Name, c.Name),
				Desc:     "container no cpu requests",
				Score:    1,
				Weight:   r.Weight,
				Proposal: "any container should set cpu requests",
			}
		}
	}
}
