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
	Param() CreateParam
	StartDiagnose(ctx context.Context) chan *Result
}

type CreateParam struct {
	Name   string
	Score  int
	Weight int
	Cli    kubernetes.Interface
}

type Creator func(d *CreateParam) Diagnostic

var Creators = map[string]Creator{}

func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
