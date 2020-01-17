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
	"sync"

	"tkestack.io/kube-jarvis/pkg/logger"
)

// Server is a http server of kube-jarvis
// all plugins can register standard APIs or extended APIs
type Server struct {
	handlers       map[string]func(http.ResponseWriter, *http.Request)
	handlersLock   sync.Mutex
	listenAndServe func(addr string, handler http.Handler) error
	handFunc       func(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// NewServer create a new Server with default values
func NewServer() *Server {
	return &Server{
		handlers:       map[string]func(http.ResponseWriter, *http.Request){},
		handlersLock:   sync.Mutex{},
		listenAndServe: http.ListenAndServe,
		handFunc:       http.HandleFunc,
	}
}

// Default is the default http server
var Default = NewServer()

// Start try setup a http server if any handler registered
func (s *Server) Start(logger logger.Logger, addr string) {
	if len(s.handlers) == 0 {
		return
	}

	if addr == "" {
		addr = ":9005"
	}

	logger.Infof("http server start at %s", addr)
	if err := s.listenAndServe(addr, nil); err != nil {
		panic(err.Error())
	}
}

// HandleFunc registered a handler for a certain path
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handlersLock.Lock()
	defer s.handlersLock.Unlock()
	if _, exist := s.handlers[pattern]; exist {
		return
	}

	s.handlers[pattern] = handler
	s.handFunc(pattern, handler)
}
