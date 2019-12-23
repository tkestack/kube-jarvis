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
package compexplorer

import (
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
)

type Auto struct {
	logger      logger.Logger
	Type        string
	Name        string
	Namespace   string
	Nodes       []string
	MasterNodes bool
	Labels      map[string]string
	exps        []Explorer
}

// NewAuto create a ComponentConfig with default value
func NewAuto(defName string, masterNodes bool) *Auto {
	return &Auto{
		Type:        TypeAuto,
		Name:        defName,
		Namespace:   "kube-system",
		Nodes:       []string{},
		MasterNodes: masterNodes,
	}
}

func (a *Auto) Complete() {
	if a.Type == "" {
		a.Type = TypeAuto
	}

	if a.Namespace == "" {
		a.Namespace = "kube-system"
	}
}

// Init do init
func (a *Auto) Init(logger logger.Logger, cli kubernetes.Interface, nodeExecutor nodeexec.Executor) error {
	specialNodes := false
	a.logger = logger
	if a.MasterNodes == true || len(a.Nodes) != 0 {
		specialNodes = true
	}

	if err := a.initNodes(cli); err != nil {
		return err
	}

	if a.Type == TypeAuto || a.Type == TypeLabel {
		a.exps = append(a.exps, NewLabelExp(a.logger, cli, a.Namespace, a.Name, a.Labels, nodeExecutor))
	}

	// only special nodes is supported for use static pod
	if (a.Type == TypeAuto || a.Type == TypeStaticPod) && specialNodes {
		a.exps = append(a.exps, NewStaticPods(a.logger, cli, a.Namespace, a.Name, a.Nodes, nodeExecutor))
	}

	if a.Type == TypeAuto || a.Type == TypeBare {
		a.exps = append(a.exps, NewBare(a.logger, a.Name, a.Nodes, nodeExecutor))
	}

	return nil
}

func (a *Auto) initNodes(cli kubernetes.Interface) error {
	// get masters as nodes
	if a.MasterNodes {
		label := labels.NewSelector()
		req, err := labels.NewRequirement("node-role.kubernetes.io/master", selection.Exists, nil)
		if err != nil {
			return errors.Wrap(err, "create master selector label failed")
		}
		label = label.Add(*req)

		masters, err := cli.CoreV1().Nodes().List(v1.ListOptions{
			LabelSelector: label.String(),
		})
		if err != nil {
			return errors.Wrapf(err, "get masters failed")
		}

		for _, n := range masters.Items {
			a.Nodes = append(a.Nodes, n.Name)
		}
	}

	// get all nodes as Auto.Nodes
	if len(a.Nodes) == 0 {
		nodes, err := cli.CoreV1().Nodes().List(v1.ListOptions{})
		if err != nil {
			return errors.Wrapf(err, "get masters failed")
		}

		for _, n := range nodes.Items {
			a.Nodes = append(a.Nodes, n.Name)
		}
	}

	return nil
}

// Component get cluster components
func (a *Auto) Component() ([]cluster.Component, error) {
	for _, exp := range a.exps {
		ok, result, err := a.tryExplore(exp)
		if err != nil {
			return nil, err
		}

		if ok {
			return result, nil
		}
	}

	return []cluster.Component{}, nil
}

func (a *Auto) tryExplore(exp Explorer) (bool, []cluster.Component, error) {
	result, err := exp.Component()
	if err != nil {
		return false, nil, errors.Wrapf(err, "component do explore failed ")
	}

	if len(result) == 0 {
		return false, result, nil
	}

	for _, c := range result {
		if c.IsRunning {
			return true, result, nil
		}
	}

	return false, result, nil
}
