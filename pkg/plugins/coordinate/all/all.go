package all

import (
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate/basic"
)

func init() {
	coordinate.Add("default", basic.NewCoordinator)
}
