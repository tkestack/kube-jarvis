package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/configmap"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/stdout"
)

func init() {
	export.Add("stdout", export.Factory{
		Creator: stdout.NewExporter,
	})
	export.Add("config-map", export.Factory{
		Creator: configmap.NewExporter,
	})
}
