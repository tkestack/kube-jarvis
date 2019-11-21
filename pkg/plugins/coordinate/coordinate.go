package coordinate

import (
	"context"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
)

// Coordinator knows how to coordinate diagnostics,exporters,evaluators
type Coordinator interface {
	// AddDiagnostic add a diagnostic to Coordinator
	AddDiagnostic(dia diagnose.Diagnostic)
	// AddExporter add a Exporter to Coordinator
	AddExporter(exporter export.Exporter)
	// AddEvaluate add a evaluate to Coordinator
	AddEvaluate(evaluate evaluate.Evaluator)
	// Run will do all diagnostics, evaluations, then export it by exporters
	Run(ctx context.Context)
}

// Creator is a factory to create a Coordinator
type Creator func(logger logger.Logger) Coordinator

// Creators store all registered Coordinator Creator
var Creators = map[string]Creator{}

// Add register a Coordinator Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
