package basic

import (
	"context"
	logger2 "github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster/fake"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/other/example"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/stdout"
	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
	"testing"
)

func TestNewDefault(t *testing.T) {
	logger := logger2.NewLogger()
	ctx := context.Background()
	d := NewCoordinator(logger, fake.NewCluster())

	d.AddDiagnostic(example.NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
		},
	}))
	d.AddExporter(stdout.NewExporter(&export.MetaData{}))
	d.Run(ctx)
}
