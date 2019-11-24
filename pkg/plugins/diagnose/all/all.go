package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/example"
)

func init() {
	diagnose.Add("example", diagnose.Factory{
		Creator: example.NewDiagnostic,
	})

	diagnose.Add("requests-limits", diagnose.Factory{
		Creator: example.NewDiagnostic,
	})
}
