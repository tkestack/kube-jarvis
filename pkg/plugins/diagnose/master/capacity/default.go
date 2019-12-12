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
package capacity

import "k8s.io/apimachinery/pkg/api/resource"

var (
	DefCapacities = []Capacity{
		{
			MaxNodeTotal: 5,
			Memory:       resource.MustParse("8000000Ki"),
			Cpu:          resource.MustParse("4000m"),
		},
		{
			MaxNodeTotal: 20,
			Memory:       resource.MustParse("16000000Ki"),
			Cpu:          resource.MustParse("4000m"),
		},
		{
			MaxNodeTotal: 100,
			Memory:       resource.MustParse("32000000Ki"),
			Cpu:          resource.MustParse("8000m"),
		},
		{
			MaxNodeTotal: 200,
			Memory:       resource.MustParse("64000000Ki"),
			Cpu:          resource.MustParse("16000m"),
		},
		{
			MaxNodeTotal: 100000,
			Memory:       resource.MustParse("128000000Ki"),
			Cpu:          resource.MustParse("16000m"),
		},
	}
)
