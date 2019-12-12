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
package diagnose

import (
	"context"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/translate"
)

// HealthyLevel means the healthy level of diagnostic result
type HealthyLevel string

const (
	// HealthyLevelFailed means diagnostic is failed
	// this is not a error but an acceptable failed
	HealthyLevelFailed HealthyLevel = "failed"
	// HealthyLevelGood means no health problem
	HealthyLevelGood HealthyLevel = "good"
	// HealthyLevelWarn means warn unHealthy
	HealthyLevelWarn HealthyLevel = "warn"
	// HealthyLevelRisk means risk unHealthy
	HealthyLevelRisk HealthyLevel = "risk"
	// HealthyLevelSerious means serious unHealthy
	HealthyLevelSerious HealthyLevel = "serious"
)

// Catalogue is the default catalogue type of a Diagnostic
type Catalogue []string

var (
	// CatalogueMaster Diagnostic diagnose controller panel status
	// master nodes status should belong to this catalogue
	CatalogueMaster Catalogue = []string{"master"}
	// CatalogueNode Diagnostics diagnose nodes status
	CatalogueNode Catalogue = []string{"node"}
	// CatalogueResource Diagnostics diagnose cluster resources status
	CatalogueResource Catalogue = []string{"resource"}
	// CatalogueOther Diagnostics have no certain catalogue
	CatalogueOther Catalogue = []string{"other"}
)

// MetaData contains core attributes of a Diagnostic
type MetaData struct {
	plugins.CommonMetaData
	// Catalogue is the catalogue type of the Diagnostic
	Catalogue Catalogue
}

// Meta return core MetaData
// this function can be use for struct implement Diagnostic interface
func (m *MetaData) Meta() MetaData {
	return *m
}

// Result is a diagnostic result item
type Result struct {
	// Level is the healthy status
	Level HealthyLevel
	// ObjName is the name of diagnosed object
	ObjName string
	// Title is the short description of Result,that is, the title of Result
	Title translate.Message
	// Desc is the full description of Result
	Desc translate.Message
	// Proposal is the full description that show how solve the healthy problem
	Proposal translate.Message
}

// StartDiagnoseParam contains all items that StartDiagnose need
type StartDiagnoseParam struct {
	// CloudType is the cloud provider type fo cluster
	CloudType string
	// Resources contains all diagnose able resources
	Resources *cluster.Resources
}

// Diagnostic diagnose some aspects of cluster
type Diagnostic interface {
	// Meta return core MetaData
	Meta() MetaData
	// StartDiagnose return a result chan that will output results
	// [chan *Result] will pop results one by one
	// diagnostic is deemed to finish if [chan *Result] is closed
	StartDiagnose(ctx context.Context, param StartDiagnoseParam) chan *Result
}

// Factory create a new Diagnostic
type Factory struct {
	// Creator is a factory function to create Diagnostic
	Creator func(d *MetaData) Diagnostic
	// SupportedClouds indicate what cloud providers will be supported of this diagnostic
	SupportedClouds []string
	// Catalogue is the catalogue type of the Diagnostic
	Catalogue Catalogue
}

// Factories store all registered Diagnostic Creator
var Factories = map[string]Factory{}

// Add register a Diagnostic Factory
func Add(typ string, f Factory) {
	Factories[typ] = f
}
