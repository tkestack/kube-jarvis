package cluster

import (
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Cluster is the abstract of target cluster
// other plugins should get Resources from Cluster
type Cluster interface {
	// Init Instantiation for cluster, it will fetch Resources
	Init() error
	// SyncResources reFetch all resource from cluster
	SyncResources() error
	// CloudType return the cloud type of Cluster
	CloudType() string
	// Resources just return fetched resources
	Resources() *Resources
}

// Factory create a new Cluster
type Factory struct {
	// Creator is a factory function to create Cluster
	Creator func(log logger.Logger, cli kubernetes.Interface, config *rest.Config) Cluster
}

// Factories store all registered Cluster Creator
var Factories = map[string]Factory{}

// Add register a Diagnostic Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}
