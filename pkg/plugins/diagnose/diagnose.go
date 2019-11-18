package diagnose

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"k8s.io/client-go/kubernetes"
)

type HealthyLevel string

const (
	HealthyLevelPass = "pass"
	HealthyLevelWarn = "warn"
	HealthyLevelRisk = "risk"
)

// Result is a diagnostic result item
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

// Diagnostic diagnose some aspects of cluster
type Diagnostic interface {
	Param() CreateParam
	// StartDiagnose return a result chan that will output results
	StartDiagnose(ctx context.Context) chan *Result
}

// CreateParam contains core attributes of a Diagnostic
type CreateParam struct {
	Logger logger.Logger
	Name   string
	Score  int
	Weight int
	Cli    kubernetes.Interface
}

// Creator is a factory to create a Diagnostic
type Creator func(d *CreateParam) Diagnostic

// Creators store all registered Diagnostic Creator
var Creators = map[string]Creator{}

// Add register a Diagnostic Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
