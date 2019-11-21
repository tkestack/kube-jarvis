package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/stdout"
)

func init() {
	export.Add("stdout", stdout.NewExporter)
}
