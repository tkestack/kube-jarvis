package configmap

import (
	"testing"

	"tkestack.io/kube-jarvis/pkg/plugins"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

func TestNewStdout(t *testing.T) {
	cli := fake.NewSimpleClientset()
	s := NewExporter(&export.MetaData{
		CommonMetaData: plugins.CommonMetaData{},
	}).(*Exporter)
	s.Cli = cli
	export.RunExporterTest(t, s)

	cm, err := cli.CoreV1().ConfigMaps(s.Namespace).Get(s.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	data := cm.Data[s.DataKey]
	t.Logf(data)

	c := export.Collector{
		Format: s.Format,
	}
	if err := c.Unmarshal([]byte(data)); err != nil {
		t.Fatalf(err.Error())
	}

	if len(c.Diagnostics) < 1 {
		t.Fatalf("no diagnostics found")
	}

	if len(c.Diagnostics[0].Results) < 1 {
		t.Fatalf("no diagnostics results found")
	}

}
