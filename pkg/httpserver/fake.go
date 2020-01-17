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

import "net/http"

// FakeResponseWriter is a fake http.ResponseWriter for handler testing
type FakeResponseWriter struct {
	StatusCode int
	RespData   []byte
	HeaderMap  http.Header
}

// NewFakeResponseWriter return an new FakeResponseWriter
func NewFakeResponseWriter() *FakeResponseWriter {
	return &FakeResponseWriter{
		HeaderMap:  http.Header{},
		StatusCode: http.StatusOK,
	}
}

// Header return the HeaderMap of  FakeResponseWriter
func (f *FakeResponseWriter) Header() http.Header {
	return f.HeaderMap
}

// Write will append response data to FakeResponseWriter.RespData
func (f *FakeResponseWriter) Write(data []byte) (int, error) {
	f.RespData = append(f.RespData, data...)
	return len(data), nil
}

// WriteHeader set the StatusCode of FakeResponseWriter
func (f *FakeResponseWriter) WriteHeader(statusCode int) {
	f.StatusCode = statusCode
}
