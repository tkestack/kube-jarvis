package configmap

import (
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestNewStdout(t *testing.T) {
	cli := fake.NewSimpleClientset()
	s := NewExporter(&export.CreateParam{
		Cli: cli,
	}).(*Exporter)

	export.RunExporterTest(t, s)

	cm, err := cli.CoreV1().ConfigMaps(s.Namespace).Get(s.Name, metav1.GetOptions{})
	if err != nil {
		t.Fatalf(err.Error())
	}

	data := cm.Data[s.DataKey]
	t.Logf(data)

	c := export.Collector{}
	if err := c.Unmarshal(s.Format, []byte(data)); err != nil {
		t.Fatalf(err.Error())
	}

	if len(c.Diagnostics) < 1 {
		t.Fatalf("no diagnostics found")
	}

	if len(c.Diagnostics[0].Results) < 1 {
		t.Fatalf("no diagnostics results found")
	}

	if len(c.Evaluations) < 1 {
		t.Fatalf("no Evaluations found")
	}

	if c.Evaluations[0].Result.Name == "" {
		t.Fatalf("empty evalution result")
	}
}
