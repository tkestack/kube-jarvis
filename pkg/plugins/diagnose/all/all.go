package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/other/example"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/resource/workload/requestslimits"
)

func init() {
	diagnose.Add("example", diagnose.Factory{
		Creator:   example.NewDiagnostic,
		Catalogue: diagnose.CatalogueOther,
	})

	diagnose.Add("requests-limits", diagnose.Factory{
		Creator:   requestslimits.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})
}
