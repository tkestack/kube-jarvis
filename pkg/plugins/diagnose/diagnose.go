package diagnose

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
)

// HealthyLevel means the healthy level of diagnostic result
type HealthyLevel string

const (
	// HealthyLevelGood means no health problem
	HealthyLevelGood = "good"
	// HealthyLevelWarn means warn unHealthy
	HealthyLevelWarn = "warn"
	// HealthyLevelRisk means risk unHealthy
	HealthyLevelRisk = "risk"
	// HealthyLevelSerious means serious unHealthy
	HealthyLevelSerious = "serious"
)

// MetaData contains core attributes of a Diagnostic
type MetaData struct {
	plugins.CommonMetaData
	// TotalScore is the score that this diagnostic has
	TotalScore float64
	// Score is th score result of this diagnostic
	// it should be set by diagnostic once diagnostic finish
	Score float64
}

// Result is a diagnostic result item
type Result struct {
	// Level is the healthy status
	Level HealthyLevel
	// Name is the short description of Result,that is, the title of Result
	Name translate.Message
	// ObjName is the name of diagnosed object
	ObjName string
	// Desc is the full description of Result
	Desc translate.Message
	// Score is the score that will be subtract from Diagnostic
	Score float64
	// Error is the error detail if diagnostic failed
	Error error
	// Proposal is the full description that show how solve the healthy problem
	Proposal translate.Message
}

// Diagnostic diagnose some aspects of cluster
type Diagnostic interface {
	// Meta return core MetaData
	Meta() MetaData
	// StartDiagnose return a result chan that will output results
	// [chan *Result] will pop results one by one
	// diagnostic is deemed to finish if [chan *Result] is closed
	StartDiagnose(ctx context.Context) chan *Result
}

// Factory create a new Diagnostic
type Factory struct {
	// Creator is a factory function to create Diagnostic
	Creator func(d *MetaData) Diagnostic
	// SupportedClouds indicate what cloud providers will be supported of this diagnostic
	SupportedClouds []string
}

// Factories store all registered Diagnostic Creator
var Factories = map[string]Factory{}

// Add register a Diagnostic Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}

// IsSupported return true if cloud type is supported by Diagnostic
func (f *Factory) IsSupported(cloud string) bool {
	if len(f.SupportedClouds) == 0 {
		return true
	}

	for _, c := range f.SupportedClouds {
		if c == cloud {
			return true
		}
	}
	return false
}
