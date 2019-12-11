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
package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	lg := NewLogger()
	lg = lg.With(map[string]string{
		"user": "test",
	}).With(map[string]string{
		"logger": "fmt",
	})

	lg.Infof("info")
	lg.Errorf("error")
	lg.Debugf("debug")
}
