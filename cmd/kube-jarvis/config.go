package main

import (
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/RayHuangCN/kube-jarvis/pkg/util"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"os"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
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
		Score     float64
		Catalogue diagnose.Catalogue
		Config    interface{}
	}

	Evaluators []struct {
		Type   string
		Name   string
		Config interface{}
	}

	Exporters []struct {
		Type   string
		Name   string
		Config interface{}
	}
}

// GetConfig return a Config struct according to content of config file
func GetConfig(file string, log logger.Logger) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "read file failed")
	}
	return getConfig(data, log)
}

func getConfig(data []byte, log logger.Logger) (*Config, error) {
	c := &Config{
		Logger: log,
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

	if err := cls.Init(); err != nil {
		return nil, errors.Wrap(err, "init cluster failed")
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

	return cr, nil
}

// GetDiagnostics create all target Diagnostics
func (c *Config) GetDiagnostics(cls cluster.Cluster, trans translate.Translator) ([]diagnose.Diagnostic, error) {
	ds := make([]diagnose.Diagnostic, 0)
	nameSet := map[string]bool{}
	for _, config := range c.Diagnostics {
		if config.Name == "" {
			config.Name = config.Type
		}

		if config.Score == 0 {
			config.Score = 100
		}

		if nameSet[config.Name] {
			return nil, fmt.Errorf("diagnostic [%s] name already exist", config.Name)
		}
		nameSet[config.Name] = true

		factory, exist := diagnose.Factories[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found diagnostic type %s", config.Type)
		}

		if !plugins.IsSupportedCloud(factory.SupportedClouds, cls.CloudType()) {
			c.Logger.Infof("diagnostic [%s] don't support cloud [%s], skipped", config.Name, cls.CloudType())
			continue
		}

		catalogue := config.Catalogue
		if catalogue == "" {
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

		if err := d.Init(); err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

// GetExporters create all target Exporters
func (c *Config) GetExporters(cls cluster.Cluster, trans translate.Translator) ([]export.Exporter, error) {
	es := make([]export.Exporter, 0)
	nameSet := map[string]bool{}
	for _, config := range c.Exporters {
		if config.Name == "" {
			config.Name = config.Type
		}

		if nameSet[config.Name] {
			return nil, fmt.Errorf("exporter [%s] name already exist", config.Name)
		}
		nameSet[config.Name] = true

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

		es = append(es, e)
	}

	return es, nil
}
