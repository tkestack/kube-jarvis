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
package main

import (
	"context"
	"flag"
	_ "tkestack.io/kube-jarvis/pkg/plugins/cluster/all"
	_ "tkestack.io/kube-jarvis/pkg/plugins/coordinate/all"
	_ "tkestack.io/kube-jarvis/pkg/plugins/diagnose/all"
	_ "tkestack.io/kube-jarvis/pkg/plugins/export/all"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "conf/default.yaml", "config file")
}

func main() {
	config, err := GetConfig(configFile)
	if err != nil {
		panic(err)
	}

	cls, err := config.GetCluster()
	if err != nil {
		panic(err)
	}

	coordinator, err := config.GetCoordinator(cls)
	if err != nil {
		panic(err)
	}

	trans, err := config.GetTranslator()
	if err != nil {
		panic(err)
	}

	diagnostics, err := config.GetDiagnostics(cls, trans)
	if err != nil {
		panic(err)
	}

	for _, d := range diagnostics {
		coordinator.AddDiagnostic(d)
	}

	exporters, err := config.GetExporters(cls, trans)
	if err != nil {
		panic(err)
	}

	for _, e := range exporters {
		coordinator.AddExporter(e)
	}

	coordinator.Run(context.Background())
}
