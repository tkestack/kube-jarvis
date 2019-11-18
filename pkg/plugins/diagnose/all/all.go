package all

import (
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/example"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose/requestslimites"
)

func init() {
	diagnose.Add("example", example.NewDiagnostic)
	diagnose.Add("requests-limits", requestslimites.NewDiagnostic)
}
