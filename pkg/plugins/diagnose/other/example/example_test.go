package example

import (
	"context"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
)

func TestDiagnostic_StartDiagnose(t *testing.T) {
	s := NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake().WithModule("diagnostics.example"),
		},
	})

	if err := s.Init(); err != nil {
		t.Fatalf(err.Error())
	}

	result := s.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{})

	for {
		res, ok := <-result
		if !ok {
			break
		}
		t.Logf("%+v", res)
	}
}
