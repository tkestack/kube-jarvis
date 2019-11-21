package example

import (
	"context"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"github.com/RayHuangCN/kube-jarvis/pkg/translate"
)

func TestDiagnostic_StartDiagnose(t *testing.T) {
	s := NewDiagnostic(&diagnose.CreateParam{
		Translator: translate.NewFake().WithModule("diagnostics.example"),
		Score:      10,
		Weight:     10,
		Cli:        nil,
	})

	result := s.StartDiagnose(context.Background())

	for {
		res, ok := <-result
		if !ok {
			break
		}
		t.Logf("%+v", res)
	}
}
