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
	"context"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
)

// Cluster is the abstract of target cluster
// other plugins should get Resources from Cluster
type Cluster interface {
	// Complete check and complete config items
	Complete() error
	// Init do Initialization for cluster, and fetching Resources
	Init(ctx context.Context, progress *plugins.Progress) error
	// CloudType return the cloud type of Cluster
	CloudType() string
	// Resources just return fetched resources
	Resources() *Resources
	// Finish will be called once diagnostic done
	Finish() error
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
