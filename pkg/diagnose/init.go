package diagnose

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type HealthyLevel string

const (
	HealthyLevelPass = "pass"
	HealthyLevelWarn = "warn"
	HealthyLevelRisk = "risk"
)

type Result struct {
	Level    HealthyLevel
	Name     string
	ObjName  string
	Desc     string
	Score    int
	Weight   int
	Error    error
	Proposal string
}

type Diagnostic interface {
	Name() string
	StartDiagnose(ctx context.Context, cli kubernetes.Interface) chan *Result
}
