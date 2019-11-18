package all

import (
	"github.com/RayHuangCN/Jarvis/pkg/plugins/export"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/export/stdout"
)

func init() {
	export.Add("stdout", stdout.NewExporter)
}
