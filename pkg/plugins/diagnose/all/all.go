package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/master/capacity"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/other/example"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/resource/workload/requestslimits"
)

func init() {
	addMasterDiagnostics()
	addResourceDiagnostics()
	addOtherDiagnostics()
}

func addMasterDiagnostics() {
	diagnose.Add(capacity.DiagnosticType, diagnose.Factory{
		Creator:   capacity.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
}

func addResourceDiagnostics() {
	diagnose.Add(requestslimits.DiagnosticType, diagnose.Factory{
		Creator:   requestslimits.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})
}

func addOtherDiagnostics() {
	diagnose.Add(example.DiagnosticType, diagnose.Factory{
		Creator:   example.NewDiagnostic,
		Catalogue: diagnose.CatalogueOther,
	})
}
