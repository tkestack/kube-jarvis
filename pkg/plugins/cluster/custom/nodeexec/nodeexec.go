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
	// UnKnowTypeErr will be returned if target node executor is not found
	UnKnowTypeErr = fmt.Errorf("unknow node executor type")
	// NoneExecutor will be returned if node executor type is "none"
	NoneExecutor = fmt.Errorf("none executor")
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
	// Type is the node executor type
	Type string
	// Namespace is the namespace to install node agent if node executor type is "agent"
	Namespace string
	// DaemonSet is the DaemonSet name if node executor type is "agent"
	DaemonSet string
	// Image is the image that will be use to create node agent DaemonSet
	Image string
	// AutoCreate indicate whether to create node agent DaemonSet if it is not exist
	// if AutoCreate is true, agent will be deleted once cluster diagnostic done
	AutoCreate bool
}

// NewConfig return a Config with default value
func NewConfig() *Config {
	return &Config{}
}

// Complete check and complete config fields
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

	if c.Image == "" {
		c.Image = "raylhuang110/kube-jarvis-agent:latest"
	}
}

// Executor return the appropriate node executor according to config value
func (c *Config) Executor(logger logger.Logger,
	cli kubernetes.Interface, config *restclient.Config) (Executor, error) {
	switch c.Type {
	case "proxy":
		return NewDaemonSetProxy(logger, cli, config, c.Namespace, c.DaemonSet, c.Image, c.AutoCreate)
	case "none":
		return nil, NoneExecutor
	}
	return nil, UnKnowTypeErr
}
