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
package nodeexec

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"tkestack.io/kube-jarvis/pkg/logger"
)

var (
	UnKnowTypeErr = fmt.Errorf("unknow node executor type")
	NoneExecutor  = fmt.Errorf("none executor")
)

// Executor get machine information
type Executor interface {
	// DoCmd do cmd on node and return output
	DoCmd(nodeName string, cmd []string) (string, string, error)
	// Finish will be called once this Executor work done
	Finish() error
}

// Config is the config of node executor
type Config struct {
	Type      string
	Namespace string
	DaemonSet string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Complete() {
	if c.Type == "" {
		c.Type = "proxy"
	}

	if c.Namespace == "" {
		c.Namespace = "kube-jarvis"
	}

	if c.DaemonSet == "" {
		c.DaemonSet = "kube-jarvis-agent"
	}
}

func (c *Config) Executor(logger logger.Logger, cli kubernetes.Interface, config *restclient.Config) (Executor, error) {
	switch c.Type {
	case "proxy":
		return NewDaemonSetProxy(logger, cli, config, c.Namespace, c.DaemonSet)
	case "none":
		return nil, NoneExecutor
	}
	return nil, UnKnowTypeErr
}
