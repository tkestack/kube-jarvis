package all

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster/custom"
)

func init() {
	cluster.Add(custom.Type, cluster.Factory{Creator: custom.NewCluster})
}
