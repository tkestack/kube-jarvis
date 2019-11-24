package configmap

import (
	"context"
	"fmt"
	"time"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Exporter save result to K8S ConfigMap with format "json" or  "yaml"
type Exporter struct {
	*export.MetaData
	export.Collector
	Namespace string
	Name      string
	Format    string
	DataKey   string
}

// NewExporter return a config-map Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	return &Exporter{
		MetaData: m,
	}
}

// Meta return core attributes
func (e *Exporter) Meta() export.MetaData {
	return *e.MetaData
}

// CoordinateFinish export save collected data to config-map
func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	e.initDefault()
	cm, err := e.getConfigMap()
	if err != nil {
		return err
	}

	data, err := e.Marshal(e.Format)
	if err != nil {
		return fmt.Errorf("unmaral data failed : %v", err)
	}

	cm.Data[e.DataKey] = string(data)
	_, err = e.Cli.CoreV1().ConfigMaps(e.Namespace).Update(cm)
	return err
}

func (e *Exporter) getConfigMap() (*v1.ConfigMap, error) {
	cm, err := e.Cli.CoreV1().ConfigMaps(e.Namespace).Get(e.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			newCm := &v1.ConfigMap{}
			newCm.Name = e.Name
			newCm.Namespace = e.Namespace
			cm, err = e.Cli.CoreV1().ConfigMaps(e.Namespace).Create(newCm)
			if err != nil {
				return nil, fmt.Errorf("create config map failed : %v", err)
			}
		}
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}

	return cm, nil
}

func (e *Exporter) initDefault() {
	if e.Name == "" {
		e.Name = "kube-jarvis"
	}

	if e.Format == "" {
		e.Format = "json"
	}

	if e.Namespace == "" {
		e.Namespace = "default"
	}

	e.DataKey = time.Now().Format("2006-01-02T15-04-05")
}
