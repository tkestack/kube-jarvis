package diagnose

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type SampleDiagnostic struct {
	result   chan *Result
	Weight   int
	NameDesc string
}

func NewSampleDiagnostic() *SampleDiagnostic {
	return &SampleDiagnostic{
		result:   make(chan *Result, 1000),
		Weight:   1,
		NameDesc: "SampleDiagnostic",
	}
}

func (s *SampleDiagnostic) Name() string {
	return s.NameDesc
}

func (s *SampleDiagnostic) StartDiagnose(ctx context.Context, cli kubernetes.Interface) chan *Result {
	go func() {
		defer close(s.result)
		s.result <- &Result{
			Level:    HealthyLevelRisk,
			Name:     "sample",
			ObjName:  "sample-obj",
			Desc:     "this is sample Diagnostic",
			Score:    10,
			Weight:   100,
			Proposal: "sample proposal",
		}
	}()
	return s.result
}
