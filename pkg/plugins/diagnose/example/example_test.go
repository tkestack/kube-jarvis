package example

import (
	"context"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"
)

func TestDiagnostic_StartDiagnose(t *testing.T) {
	s := NewDiagnostic(&diagnose.CreateParam{
		Score:  10,
		Weight: 10,
		Cli:    nil,
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
