package coordinate

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"
)

// Coordinator knows how to coordinate diagnostics,exporters,evaluators
type Coordinator interface {
	AddDiagnostic(dia diagnose.Diagnostic)
	AddExporter(exporter export.Exporter)
	AddEvaluate(evaluate evaluate.Evaluator)
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
