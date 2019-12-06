package compexplorer

import "tkestack.io/kube-jarvis/pkg/plugins/cluster"

const (
	TypeStaticPod = "StaticPod"
	TypeLabel     = "Label"
	TypeBare      = "Bare"
)

// Explorer get component information
type Explorer interface {
	// Component get cluster components
	Component() ([]cluster.Component, error)
}
