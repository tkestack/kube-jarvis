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
package all

import (
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/args/apiserver"
	controller_manager "tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/args/controller-manager"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/args/etcd"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/args/scheduler"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/capacity"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/master/components"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/node/ha"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/node/iptables"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/node/status"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/node/sys"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/other/example"
	hpaip "tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/hpa/ip"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/affinity"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/batch"
	workloadha "tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/ha"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/healthcheck"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/pdb"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/requestslimits"
	workloadStatus "tkestack.io/kube-jarvis/pkg/plugins/diagnose/resource/workload/status"
)

func init() {
	addMasterDiagnostics()
	addResourceDiagnostics()
	addOtherDiagnostics()
	addNodeDiagnostics()
	addNodeStatusDiagnostics()
	addMasterArgDiagnostics()
}

func addMasterDiagnostics() {
	diagnose.Add(capacity.DiagnosticType, diagnose.Factory{
		Creator:   capacity.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})

	diagnose.Add(components.DiagnosticType, diagnose.Factory{
		Creator:   components.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
}

func addResourceDiagnostics() {
	diagnose.Add(requestslimits.DiagnosticType, diagnose.Factory{
		Creator:   requestslimits.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(hpaip.DiagnosticType, diagnose.Factory{
		Creator:   hpaip.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(healthcheck.DiagnosticType, diagnose.Factory{
		Creator:   healthcheck.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(affinity.DiagnosticType, diagnose.Factory{
		Creator:   affinity.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(pdb.DiagnosticType, diagnose.Factory{
		Creator:   pdb.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(batch.DiagnosticType, diagnose.Factory{
		Creator:   batch.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(workloadha.DiagnosticType, diagnose.Factory{
		Creator:   workloadha.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})

	diagnose.Add(workloadStatus.DiagnosticType, diagnose.Factory{
		Creator:   workloadStatus.NewDiagnostic,
		Catalogue: diagnose.CatalogueResource,
	})
}

func addOtherDiagnostics() {
	diagnose.Add(example.DiagnosticType, diagnose.Factory{
		Creator:   example.NewDiagnostic,
		Catalogue: diagnose.CatalogueOther,
	})
}

func addNodeDiagnostics() {
	diagnose.Add(sys.DiagnosticType, diagnose.Factory{
		Creator:   sys.NewDiagnostic,
		Catalogue: diagnose.CatalogueNode,
	})
	diagnose.Add(iptables.DiagnosticType, diagnose.Factory{
		Creator:   iptables.NewDiagnostic,
		Catalogue: diagnose.CatalogueNode,
	})
	diagnose.Add(ha.DiagnosticType, diagnose.Factory{
		Creator:   ha.NewDiagnostic,
		Catalogue: diagnose.CatalogueNode,
	})
}

func addNodeStatusDiagnostics() {
	diagnose.Add(status.DiagnosticType, diagnose.Factory{
		Creator:   status.NewDiagnostic,
		Catalogue: diagnose.CatalogueNode,
	})
}

func addMasterArgDiagnostics() {
	diagnose.Add(apiserver.DiagnosticType, diagnose.Factory{
		Creator:   apiserver.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
	diagnose.Add(scheduler.DiagnosticType, diagnose.Factory{
		Creator:   scheduler.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
	diagnose.Add(controller_manager.DiagnosticType, diagnose.Factory{
		Creator:   controller_manager.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
	diagnose.Add(etcd.DiagnosticType, diagnose.Factory{
		Creator:   etcd.NewDiagnostic,
		Catalogue: diagnose.CatalogueMaster,
	})
}
