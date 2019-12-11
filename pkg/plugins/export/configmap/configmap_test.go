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
package configmap

import (
	"testing"

	"tkestack.io/kube-jarvis/pkg/plugins"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

func TestNewStdout(t *testing.T) {
	cli := fake.NewSimpleClientset()
	s := NewExporter(&export.MetaData{
		CommonMetaData: plugins.CommonMetaData{},
	}).(*Exporter)
	s.Cli = cli
	export.RunExporterTest(t, s)

	cm, err := cli.CoreV1().ConfigMaps(s.Namespace).Get(s.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	data := cm.Data[s.DataKey]
	t.Logf(data)

	c := export.Collector{
		Format: s.Format,
	}
	if err := c.Unmarshal([]byte(data)); err != nil {
		t.Fatalf(err.Error())
	}

	if len(c.Diagnostics) < 1 {
		t.Fatalf("no diagnostics found")
	}

	if len(c.Diagnostics[0].Results) < 1 {
		t.Fatalf("no diagnostics results found")
	}

}
