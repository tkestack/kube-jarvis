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
package plugins

import (
	"fmt"
	"testing"
)

func TestIsSupportedCloud(t *testing.T) {
	var cases = []struct {
		supported bool
		clouds    []string
		cloud     string
	}{
		{
			supported: true,
			clouds:    []string{},
			cloud:     "123",
		},
		{
			supported: true,
			clouds: []string{
				"123",
			},
			cloud: "123",
		},
		{
			supported: false,
			clouds: []string{
				"321",
			},
			cloud: "123",
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			if IsSupportedCloud(cs.clouds, cs.cloud) != cs.supported {
				t.Fatalf("shoud %v", cs.supported)
			}
		})
	}
}
