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
package httpserver

import (
	"net/http"
	"testing"
)

func TestFakeResponseWriter_Header(t *testing.T) {
	f := NewFakeResponseWriter()
	f.HeaderMap.Set("a", "b")
	h := f.Header()
	if h.Get("a") != "b" {
		t.Fatalf("want a but get %s", h.Get("a"))
	}
}

func TestFakeResponseWriter_Write(t *testing.T) {
	f := NewFakeResponseWriter()
	str := "test"
	n, err := f.Write([]byte(str))
	if err != nil {
		t.Fatalf(err.Error())
	}

	if n != len(str) {
		t.Fatalf("want %d, but get %d", len(str), n)
	}

	if string(f.RespData) != str {
		t.Fatalf("want %s, but get %s", str, string(f.RespData))
	}
}

func TestFakeResponseWriter_WriteHeader(t *testing.T) {
	f := NewFakeResponseWriter()
	f.WriteHeader(http.StatusOK)
	if f.StatusCode != http.StatusOK {
		t.Fatalf("want 200 but get %d", f.StatusCode)
	}
}
