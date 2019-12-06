package all

import (
	"tkestack.io/kube-jarvis/pkg/plugins/export"
	"tkestack.io/kube-jarvis/pkg/plugins/export/file"
	"tkestack.io/kube-jarvis/pkg/plugins/export/stdout"
)

func init() {
	export.Add(stdout.ExporterType, export.Factory{
		Creator: stdout.NewExporter,
	})

	export.Add(file.ExporterType, export.Factory{
		Creator: file.NewExporter,
	})
	/*
		export.Add(configmap.ExporterType, export.Factory{
			Creator: configmap.NewExporter,
		})*/
}
