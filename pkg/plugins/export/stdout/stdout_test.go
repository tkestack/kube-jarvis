package stdout

import (
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/export"
)

func TestNewStdout(t *testing.T) {
	s := NewExporter(&export.MetaData{}).(*Exporter)
	export.RunExporterTest(t, s)
}
