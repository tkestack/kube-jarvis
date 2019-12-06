package example

import (
	"context"
	"testing"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
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
