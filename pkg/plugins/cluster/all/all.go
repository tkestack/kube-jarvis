package all

import (
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom"
)

func init() {
	cluster.Add(custom.Type, cluster.Factory{Creator: custom.NewCluster})
}
