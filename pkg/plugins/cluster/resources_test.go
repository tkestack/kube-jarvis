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

package cluster

import (
	"fmt"
	"testing"
)

func TestResourcesFilter_Compile(t *testing.T) {
	cases := []struct {
		filter ResourcesFilter
		pass   bool
	}{
		{
			filter: ResourcesFilter{
				{
					Namespace: "**",
				},
			},
			pass: false,
		},
		{
			filter: ResourcesFilter{
				{
					Namespace: ".*",
				},
			},
			pass: true,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			err := cs.filter.Compile()
			if (err == nil) != cs.pass {
				t.Fatalf("want pass = %v, but not", cs.pass)
			}
		})
	}
}

func TestResourcesFilter_Filtered(t *testing.T) {
	cases := []struct {
		filter   ResourcesFilter
		ns       string
		kind     string
		name     string
		filtered bool
	}{
		{
			filter: ResourcesFilter{
				{
					Namespace: "test",
				},
			},
			ns:       "test",
			filtered: true,
		},
		{
			filter: ResourcesFilter{
				{
					Namespace: "test",
					Kind:      "Pod",
					Name:      "Pod1",
				},
			},
			ns:       "test1",
			kind:     "Pod2",
			name:     "Pod2",
			filtered: false,
		},
		{
			ns:       "test",
			filtered: false,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			if err := cs.filter.Compile(); err != nil {
				t.Fatalf(err.Error())
			}

			if cs.filtered != cs.filter.Filtered(cs.ns, cs.kind, cs.name) {
				t.Fatalf("want %v but not", cs.filtered)
			}
		})
	}
}
