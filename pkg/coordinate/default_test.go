package coordinate

import (
	"context"
	"testing"

	logger2 "github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/evaluate"
	"github.com/RayHuangCN/Jarvis/pkg/export"
)

func TestNewDefault(t *testing.T) {
	logger := logger2.NewLogger()
	ctx := context.Background()
	d := NewDefault(logger, nil)

	d.AddDiagnostic(diagnose.NewSampleDiagnostic())
	d.AddEvaluate(evaluate.NewSumEva())
	d.AddExporter(export.NewStdout())
	d.Run(ctx)
}
