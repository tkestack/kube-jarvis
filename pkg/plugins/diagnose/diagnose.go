package diagnose

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
	"k8s.io/client-go/kubernetes"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
)

// HealthyLevel means the healthy level of diagnostic result
type HealthyLevel string

const (
	// HealthyLevelWarn means warn unHealthy
	HealthyLevelWarn = "warn"
	// HealthyLevelRisk means risk unHealthy
	HealthyLevelRisk = "risk"
	// HealthyLevelRisk means serious unHealthy
	HealthyLevelSerious = "serious"
)

// Result is a diagnostic result item
type Result struct {
	Level    HealthyLevel
	Name     translate.Message
	ObjName  string
	Desc     translate.Message
	Score    int
	Weight   int
	Error    error
	Proposal translate.Message
}

// Diagnostic diagnose some aspects of cluster
type Diagnostic interface {
	// Param return core attributes
	Param() CreateParam
	// StartDiagnose return a result chan that will output results
	StartDiagnose(ctx context.Context) chan *Result
}

// CreateParam contains core attributes of a Diagnostic
type CreateParam struct {
	Cli        kubernetes.Interface
	Translator translate.Translator
	Logger     logger.Logger
	Type       string
	Name       string
	Score      int
	Weight     int
	CloudType  string
}

// Creator is a factory to create a Diagnostic
type Creator func(d *CreateParam) Diagnostic

// Creators store all registered Diagnostic Creator
var Creators = map[string]Creator{}

// Add register a Diagnostic Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
