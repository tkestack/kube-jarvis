/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package cluster

import (
	ar "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// Machine is the contains low level system information of a node
type Machine struct {
	// SysCtl is the OS system param from command "sysctl -a"
	SysCtl map[string]string
}

// ComponentName it the type of a component like
type ComponentName string

// core components of cluster
const (
	ComponentApiserver         = "kube-apiserver"
	ComponentScheduler         = "kube-scheduler"
	ComponentControllerManager = "kube-controller-manager"
	ComponentETCD              = "etcd"
	ComponentKubeProxy         = "kube-proxy"
	ComponentCoreDNS           = "coredns"
	ComponentKubeDNS           = "kube-dns"
	ComponentKubelet           = "kubelet"
	ComponentDockerd           = "dockerd"
	ComponentContainerd        = "containerd"
)

// Component is the com￿mon data of a component like kube-apiserver, etcd, schedule....
type Component struct {
	// Name is the full name of the component
	Name string
	// Node is the node name that this component run at
	Node string
	// Args is the command line of the component
	Args map[string]string
	// IsRunning is true if Component run normally
	IsRunning bool
	// Error if not nil if fetching this Component failed
	Error error
	// Pod is not nil if this Component run as pod
	Pod *corev1.Pod
}

type Resources struct {
	Nodes                           *corev1.NodeList
	PersistentVolumes               *corev1.PersistentVolumeList
	ComponentStatuses               *corev1.ComponentStatusList
	SystemNamespace                 *corev1.Namespace
	Pods                            *corev1.PodList
	PodTemplates                    *corev1.PodTemplateList
	PersistentVolumeClaims          *corev1.PersistentVolumeClaimList
	ConfigMaps                      *corev1.ConfigMapList
	Services                        *corev1.ServiceList
	Secrets                         *corev1.SecretList
	ServiceAccounts                 *corev1.ServiceAccountList
	ResourceQuotas                  *corev1.ResourceQuotaList
	LimitRanges                     *corev1.LimitRangeList
	MutatingWebhookConfigurations   *ar.MutatingWebhookConfigurationList
	ValidatingWebhookConfigurations *ar.ValidatingWebhookConfigurationList
	Namespaces                      *corev1.NamespaceList

	CoreComponents map[string][]Component
	Machines       map[string]Machine
	// Extra is a CloudType special resources
	Extra interface{}
}

// NewResources return a new Resources
func NewResources() *Resources {
	return &Resources{
		CoreComponents: map[string][]Component{},
		Machines:       map[string]Machine{},
	}
}
