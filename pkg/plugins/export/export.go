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
package export

import (
	"context"
	"tkestack.io/kube-jarvis/pkg/plugins"
)

// MetaData contains core attributes of a Exporter
type MetaData struct {
	plugins.MetaData
}

// Meta return core MetaData
// this function can be use for struct implement Exporter interface
func (m *MetaData) Meta() MetaData {
	return *m
}

// Exporter export all steps and results with special way or special format
type Exporter interface {
	// Complete check and complete config items
	Complete() error
	// Meta return core attributes
	Meta() MetaData
	// Export export result
	Export(ctx context.Context, result *AllResult) error
}

// Factory create a new Exporter
type Factory struct {
	// Creator is a factory function to create Exporter
	Creator func(d *MetaData) Exporter
	// SupportedClouds indicate what cloud providers will be supported of this exporter
	SupportedClouds []string
}

// Factories store all registered Exporter Creator
var Factories = map[string]Factory{}

// Add register a Exporter Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}
