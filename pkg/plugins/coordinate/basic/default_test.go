package basic

import (
	"context"
	"testing"
	logger2 "tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/other/example"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
	"tkestack.io/kube-jarvis/pkg/plugins/export/stdout"
	"tkestack.io/kube-jarvis/pkg/translate"
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
