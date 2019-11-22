package basic

import (
	"context"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/stdout"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate/sum"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"

	logger2 "github.com/RayHuangCN/kube-jarvis/pkg/logger"
)

func TestNewDefault(t *testing.T) {
	logger := logger2.NewLogger()
	ctx := context.Background()
	d := NewCoordinator(logger)

	d.AddDiagnostic(example.NewDiagnostic(&diagnose.CreateParam{
		Translator: translate.NewFake(),
		Score:      10,
	}))
	d.AddEvaluate(sum.NewEvaluator(&evaluate.CreateParam{
		Translator: translate.NewFake(),
	}))
	d.AddExporter(stdout.NewExporter(&export.CreateParam{}))
	d.Run(ctx)
}
