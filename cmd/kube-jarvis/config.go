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
package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/util"

	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	"tkestack.io/kube-jarvis/pkg/plugins/export"

	"tkestack.io/kube-jarvis/pkg/logger"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
)

// Config is the struct for config file
type Config struct {
	Logger logger.Logger
	Global struct {
		Trans string
		Lang  string
	}

	Cluster struct {
		Type       string
		Kubeconfig string
		Config     interface{}
	}

	Coordinator struct {
		Type   string
		Config interface{}
	}

	Diagnostics []struct {
		Type      string
		Name      string
		Catalogue diagnose.Catalogue
		Config    interface{}
	}
	Exporters []struct {
		Type   string
		Name   string
		level  string
		Config interface{}
	}
}

// GetConfig return a Config struct according to content of config file
func GetConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read file failed")
	}
	return getConfig(data)
}

func getConfig(data []byte) (*Config, error) {
	c := &Config{
		Logger: logger.NewLogger(),
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, errors.Wrap(err, "unmarshal data failed")
	}

	return c, nil
}

// GetTranslator return a translate.Translator
func (c *Config) GetTranslator() (translate.Translator, error) {
	return translate.NewDefault(c.Global.Trans, "en", c.Global.Lang)
}

// GetCluster create a cluster.Cluster
func (c *Config) GetCluster() (cluster.Cluster, error) {
	config, err := clientcmd.BuildConfigFromFlags("", c.Cluster.Kubeconfig)
	if err != nil {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err.Error())
		}

		config, err = clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s/.kube/config", home))
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("failed to create client-go client:" + err.Error())
	}

	factory, exist := cluster.Factories[c.Cluster.Type]
	if !exist {
		return nil, fmt.Errorf("can not found cluster type %s", c.Cluster.Type)
	}

	cls := factory.Creator(c.Logger.With(map[string]string{
		"cluster": c.Cluster.Type,
	}), clientset, config)

	if err := util.InitObjViaYaml(cls, c.Cluster.Config); err != nil {
		return nil, errors.Wrap(err, "init cluster config failed")
	}

	if err := cls.Complete(); err != nil {
		return nil, err
	}
	return cls, nil
}

// GetCoordinator return create a coordinate.Coordinator
func (c *Config) GetCoordinator(cls cluster.Cluster) (coordinate.Coordinator, error) {
	if c.Coordinator.Type == "" {
		c.Coordinator.Type = "default"
	}

	creator, exist := coordinate.Creators[c.Coordinator.Type]
	if !exist {
		return nil, fmt.Errorf("can not found coordinate type %s", c.Coordinator.Type)
	}

	cr := creator(c.Logger.With(map[string]string{
		"coordinator": c.Coordinator.Type,
	}), cls)

	if err := util.InitObjViaYaml(cr, c.Coordinator.Config); err != nil {
		return nil, err
	}

	if err := cr.Complete(); err != nil {
		return nil, err
	}

	return cr, nil
}

// GetDiagnostics create all target Diagnostics
func (c *Config) GetDiagnostics(cls cluster.Cluster, trans translate.Translator) ([]diagnose.Diagnostic, error) {
	ds := make([]diagnose.Diagnostic, 0)
	for _, config := range c.Diagnostics {
		factory, exist := diagnose.Factories[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found diagnostic type %s", config.Type)
		}

		if !plugins.IsSupportedCloud(factory.SupportedClouds, cls.CloudType()) {
			c.Logger.Infof("diagnostic [%s] don't support cloud [%s], skipped", config.Name, cls.CloudType())
			continue
		}

		catalogue := config.Catalogue
		if len(catalogue) == 0 {
			catalogue = factory.Catalogue
		}

		d := factory.Creator(&diagnose.MetaData{
			CommonMetaData: plugins.CommonMetaData{
				Translator: trans.WithModule("diagnostics." + config.Type),
				Logger: c.Logger.With(map[string]string{
					"diagnostic": config.Name,
				}),
				Type: config.Type,
				Name: config.Name,
			},
			Catalogue: catalogue,
		})

		if err := util.InitObjViaYaml(d, config.Config); err != nil {
			return nil, err
		}

		if err := d.Complete(); err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

// GetExporters create all target Exporters
func (c *Config) GetExporters(cls cluster.Cluster, trans translate.Translator) ([]export.Exporter, error) {
	es := make([]export.Exporter, 0)
	for _, config := range c.Exporters {
		factory, exist := export.Factories[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found exporter type %s", config.Type)
		}

		if !plugins.IsSupportedCloud(factory.SupportedClouds, cls.CloudType()) {
			c.Logger.Infof("diagnostic [%s] don't support cloud [%s], skipped", config.Name, cls.CloudType())
			continue
		}

		e := factory.Creator(&export.MetaData{
			CommonMetaData: plugins.CommonMetaData{
				Translator: trans.WithModule("diagnostics." + config.Type),
				Logger: c.Logger.With(map[string]string{
					"diagnostic": config.Name,
				}),
				Type: config.Type,
				Name: config.Name,
			},
		})

		if err := util.InitObjViaYaml(e, config.Config); err != nil {
			return nil, err
		}

		if err := e.Complete(); err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return es, nil
}
