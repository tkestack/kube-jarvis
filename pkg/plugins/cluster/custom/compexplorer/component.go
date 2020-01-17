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
package compexplorer

import (
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
)

const (
	// TypeStaticPod is the type of StaticPod explore
	TypeStaticPod = "StaticPod"
	// TypeLabel is the type of Label explore
	TypeLabel = "Label"
	// TypeBare is the type of Bare explore
	TypeBare = "Bare"
	// TypeAuto is the type of Auto explore
	TypeAuto = "Auto"
)

// Explorer get component information
type Explorer interface {
	// Component get cluster components
	Component() ([]cluster.Component, error)
	// Finish will be called once every thing done
	Finish() error
}
