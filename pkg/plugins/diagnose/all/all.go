package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/requestslimits"
)

func init() {
	diagnose.Add("example", example.NewDiagnostic)
	diagnose.Add("requests-limits", requestslimits.NewDiagnostic)
}
