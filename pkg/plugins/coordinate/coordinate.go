package coordinate

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"
)

type Coordinator interface {
	AddDiagnostic(dia diagnose.Diagnostic)
	AddExporter(exporter export.Exporter)
	AddEvaluate(evaluate evaluate.Evaluator)
	Run(ctx context.Context)
}

type Creator func(logger logger.Logger) Coordinator

var Creators = map[string]Creator{}

func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
