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
	"fmt"
	"net/http"
	"testing"

	"tkestack.io/kube-jarvis/pkg/logger"
)

func TestServer_HandleFunc(t *testing.T) {
	s := NewServer()
	s.handFunc = func(pattern string, handler func(http.ResponseWriter, *http.Request)) {}
	s.HandleFunc("test", func(writer http.ResponseWriter, request *http.Request) {})
	_, exist := s.handlers["test"]
	if !exist {
		t.Fatalf("register handler failed")
	}
}

func TestServer_Start(t *testing.T) {
	var cases = []struct {
		wantStared      bool
		registerHandler bool
	}{
		{
			wantStared:      true,
			registerHandler: true,
		},
		{
			wantStared:      false,
			registerHandler: false,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			s := NewServer()
			started := false
			s.listenAndServe = func(addr string, handler http.Handler) error {
				started = true
				return nil
			}

			if cs.registerHandler {
				s.HandleFunc("test", func(writer http.ResponseWriter, request *http.Request) {
				})
			}

			s.Start(logger.NewLogger(), "")

			if cs.wantStared != started {
				t.Fatalf("want %v but not", cs.wantStared)
			}
		})
	}
}
