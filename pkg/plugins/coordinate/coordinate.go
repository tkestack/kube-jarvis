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
package coordinate

import (
	"context"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/logger"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

// Coordinator knows how to coordinate diagnostics,exporters,evaluators
type Coordinator interface {
	// Complete check and complete config items
	Complete() error
	// AddDiagnostic add a diagnostic to Coordinator
	AddDiagnostic(dia diagnose.Diagnostic)
	// AddExporter add a Exporter to Coordinator
	AddExporter(exporter export.Exporter)
	// Run will do all diagnostics, evaluations, then export it by exporters
	Run(ctx context.Context)
}

// Creator is a factory to create a Coordinator
type Creator func(logger logger.Logger, cls cluster.Cluster) Coordinator

// Creators store all registered Coordinator Creator
var Creators = map[string]Creator{}

// Add register a Coordinator Creator
func Add(typ string, creator Creator) {
	Creators[typ] = creator
}
