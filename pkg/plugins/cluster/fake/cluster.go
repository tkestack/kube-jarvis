package fake

import "github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"

type Cluster struct {
	Res *cluster.Resources
}

func NewCluster() *Cluster {
	return &Cluster{Res: cluster.NewResources()}
}

// Init Instantiation for cluster, it will fetch Resources
func (c *Cluster) Init() error {
	return nil
}

// SyncResources fetch all resource from cluster
func (c *Cluster) SyncResources() error {
	c.Res = cluster.NewResources()
	return nil
}

// CloudType return the cloud type of Cluster
func (c *Cluster) CloudType() string {
	return "fake"
}

// Machine return the low level information of a node
func (c *Cluster) Resources() *cluster.Resources {
	return c.Res
}
