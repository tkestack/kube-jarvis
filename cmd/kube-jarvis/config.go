package main

import (
	"fmt"
	"io/ioutil"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config is the struct for config file
type Config struct {
	Logger logger.Logger
	Global struct {
		Trans   string
		Lang    string
		Cluster struct {
			Kubeconfig string
		}
	}

	Coordinate struct {
		Type   string
		Config interface{}
	}

	Diagnostics []struct {
		Type   string
		Name   string
		Score  int
		Weight int
		Config interface{}
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

// GetClusterClient create a k8s client
func (c *Config) GetClusterClient() (kubernetes.Interface, error) {
	if c.Global.Cluster.Kubeconfig == "fake" {
		return fake.NewSimpleClientset(), nil
	}

	if c.Global.Cluster.Kubeconfig == "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, errors.Wrap(err, "inCluster config failed")
		}
		return kubernetes.NewForConfig(config)
	}

	config, err := clientcmd.BuildConfigFromFlags("", c.Global.Cluster.Kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "BuildConfigFromFlags failed")
	}

	return kubernetes.NewForConfig(config)
}

// GetCoordinator return create a coordinate.Coordinator
func (c *Config) GetCoordinator() (coordinate.Coordinator, error) {
	creator, exist := coordinate.Creators[c.Coordinate.Type]
	if !exist {
		return nil, fmt.Errorf("can not found coordinate type %s", c.Coordinate.Type)
	}

	cr := creator(c.Logger.With(map[string]string{
		"coordinator": "default",
	}))
	if err := InitObjViaYaml(cr, c.Coordinate.Config); err != nil {
		return nil, err
	}

	return cr, nil
}

// GetDiagnostics create all target Diagnostics
func (c *Config) GetDiagnostics(cli kubernetes.Interface, trans translate.Translator) ([]diagnose.Diagnostic, error) {
	ds := make([]diagnose.Diagnostic, 0)
	for _, config := range c.Diagnostics {
		creator, exist := diagnose.Creators[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found diagnostic type %s", config.Type)
		}

		d := creator(&diagnose.CreateParam{
			Translator: trans.WithModule("diagnostics." + config.Type),
			Logger: c.Logger.With(map[string]string{
				"diagnostic": config.Name,
			}),
			Name:   config.Name,
			Score:  config.Score,
			Weight: config.Weight,
			Cli:    cli,
		})

		if err := InitObjViaYaml(d, config.Config); err != nil {
			return nil, err
		}

		ds = append(ds, d)
	}

	return ds, nil
}

// GetEvaluators create all target Evaluators
func (c *Config) GetEvaluators(trans translate.Translator) ([]evaluate.Evaluator, error) {
	es := make([]evaluate.Evaluator, 0)
	for _, config := range c.Evaluators {
		creator, exist := evaluate.Creators[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found evaluator type %s", config.Type)
		}

		e := creator(&evaluate.CreateParam{
			Translator: trans.WithModule("evaluators." + config.Type),
			Logger: c.Logger.With(map[string]string{
				"evaluator": config.Name,
			}),
			Name: config.Name,
		})

		if err := InitObjViaYaml(e, config.Config); err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return es, nil
}

// GetExporters create all target Exporters
func (c *Config) GetExporters() ([]export.Exporter, error) {
	es := make([]export.Exporter, 0)
	for _, config := range c.Exporters {
		creator, exist := export.Creators[config.Type]
		if !exist {
			return nil, fmt.Errorf("can not found exporter type %s", config.Type)
		}

		e := creator(&export.CreateParam{
			Logger: c.Logger.With(map[string]string{
				"exporter": config.Name,
			}),
			Name: config.Name,
		})

		if err := InitObjViaYaml(e, config.Config); err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return es, nil
}

// InitObjViaYaml marshal "config" to yaml data, then unMarshal data to "obj"
func InitObjViaYaml(obj interface{}, config interface{}) error {
	if obj == nil || config == nil {
		return nil
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, obj)
}
