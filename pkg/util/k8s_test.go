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
package util

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_Quantity(t *testing.T) {
	a := resource.MustParse("7.63Gi")
	if MemQuantityStr(&a) != "7.63GB" {
		t.Fatalf("shoud be 7.63GB")
	}

	a = resource.MustParse("1100m")
	t.Log(CpuQuantityStr(&a))
}
