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
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

// MemQuantityStr returns a string in GB
func MemQuantityStr(q *resource.Quantity) string {
	return fmt.Sprintf("%.2f", float64(q.Value())/1024/1024/1024) + "GB"
}

// CpuQuantityStr returns a string in Core
func CpuQuantityStr(q *resource.Quantity) string {
	return fmt.Sprintf("%.2f", float64(q.MilliValue())/1000) + "Core"
}
