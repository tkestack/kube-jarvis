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
	"testing"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose/other/example"
	"tkestack.io/kube-jarvis/pkg/translate"
)

// RunExporterTest is a tool function for exporter testing
func RunExporterTest(t *testing.T, exporter Exporter) {
	ctx := context.Background()
	_ = exporter.CoordinateBegin(ctx)

	d := example.NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
			Type:       "example",
		},
	})
	// Diagnostic
	fatalIfError(t, exporter.DiagnosticBegin(ctx, d))
	result, err := d.StartDiagnose(ctx, diagnose.StartDiagnoseParam{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	for {
		st, ok := <-result
		if !ok {
			break
		}

		fatalIfError(t, exporter.DiagnosticResult(ctx, d, st))
	}

	fatalIfError(t, exporter.DiagnosticFinish(ctx, d))
	// Evaluation
	fatalIfError(t, exporter.CoordinateFinish(ctx))
}

func fatalIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}
