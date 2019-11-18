package _default

import (
	"context"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/export/stdout"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate/sum"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/example"

	logger2 "github.com/RayHuangCN/Jarvis/pkg/logger"
)

func TestNewDefault(t *testing.T) {
	logger := logger2.NewLogger()
	ctx := context.Background()
	d := NewCoordinator(logger)

	d.AddDiagnostic(example.NewDiagnostic(&diagnose.CreateParam{
		Score:  10,
		Weight: 10,
	}))
	d.AddEvaluate(sum.NewEvaluator(&evaluate.CreateParam{}))
	d.AddExporter(stdout.NewExporter(&export.CreateParam{}))
	d.Run(ctx)
}
