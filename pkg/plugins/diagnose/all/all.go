package all

import (
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/example"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/requestslimits"
)

func init() {
	diagnose.Add("example", example.NewDiagnostic)
	diagnose.Add("requests-limits", requestslimits.NewDiagnostic)
}
