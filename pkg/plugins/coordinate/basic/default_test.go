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
package basic

import (
	"context"
	"testing"
	logger2 "tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/other/example"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
	"tkestack.io/kube-jarvis/pkg/plugins/export/stdout"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func TestNewDefault(t *testing.T) {
	logger := logger2.NewLogger()
	ctx := context.Background()
	d := NewCoordinator(logger, fake.NewCluster())
	_ = d.Complete()

	d.AddDiagnostic(example.NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
		},
	}))
	d.AddExporter(stdout.NewExporter(&export.MetaData{}))
	d.Run(ctx)
}
