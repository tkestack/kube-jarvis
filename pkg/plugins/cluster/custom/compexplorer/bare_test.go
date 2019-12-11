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
	"fmt"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
)

type fakeNodeExecutor struct {
	success bool
}

func (f *fakeNodeExecutor) DoCmd(nodeName string, cmd []string) (string, string, error) {
	out := `kube-apiserver
-a=123
-b=321
`
	if !f.success {
		return "", "", nil
	}
	return out, "", nil
}

func TestBare_Component(t *testing.T) {
	cases := []struct {
		success bool
	}{
		{
			success: true,
		},
		{
			success: false,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs), func(t *testing.T) {
			f := &fakeNodeExecutor{success: cs.success}
			b := NewBare(logger.NewLogger(), "kube-apiserver", []string{"node1"}, f)
			cmp, err := b.Component()
			if err != nil {
				t.Fatalf(err.Error())
			}

			if len(cmp) != 1 {
				t.Fatalf("want len 1 but get %d", len(cmp))
			}

			if cmp[0].IsRunning != cs.success {
				t.Fatalf("IsRuning wrong")
			}

			if !cs.success {
				return
			}

			if cmp[0].Args["a"] != "123" {
				t.Fatalf("want key a valuer 123 but get %s", cmp[0].Args["a"])
			}

			if cmp[0].Args["b"] != "321" {
				t.Fatalf("want key a valuer 321 but get %s", cmp[0].Args["a"])
			}

		})
	}
}
