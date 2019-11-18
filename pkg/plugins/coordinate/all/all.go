package all

import (
	"github.com/RayHuangCN/Jarvis/pkg/plugins/coordinate"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/coordinate/basic"
)

func init() {
	coordinate.Add("default", basic.NewCoordinator)
}
