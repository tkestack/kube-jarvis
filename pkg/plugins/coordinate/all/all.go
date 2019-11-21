package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate/basic"
)

func init() {
	coordinate.Add("default", basic.NewCoordinator)
}
