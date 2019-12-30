package stdout

import (
	"context"
	"testing"
	"time"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func TestExporter_Complete(t *testing.T) {
	s := Exporter{}
	if err := s.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	if s.Format != "fmt" {
		t.Fatalf("Format default value should be 'fmt'")
	}

	if s.Level != diagnose.HealthyLevelGood {
		t.Fatalf("Level default value should be 'good'")
	}
}

func TestExporter_Export(t *testing.T) {

	s := NewExporter(&export.MetaData{
		MetaData: plugins.MetaData{
			Translator: translate.NewFake(),
			Logger:     logger.NewLogger(),
			Type:       ExporterType,
		},
	})

	if err := s.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	if err := s.Export(context.Background(), &export.AllResult{
		StartTime: time.Now(),
		EndTime:   time.Now(),
		Statistics: map[diagnose.HealthyLevel]int{
			diagnose.HealthyLevelGood:    1,
			diagnose.HealthyLevelWarn:    1,
			diagnose.HealthyLevelRisk:    1,
			diagnose.HealthyLevelSerious: 1,
			diagnose.HealthyLevelFailed:  1,
		},
		Diagnostics: []*export.DiagnosticResultItem{
			{
				Results: []*diagnose.Result{
					{
						Level: diagnose.HealthyLevelGood,
						Title: "good",
					},
					{
						Level: diagnose.HealthyLevelWarn,
						Title: "warn",
					}, {
						Level: diagnose.HealthyLevelRisk,
						Title: "risk",
					},
					{
						Level: diagnose.HealthyLevelSerious,
						Title: "serious",
					},
					{
						Level: diagnose.HealthyLevelFailed,
						Title: "failed",
					},
				},
				Statistics: map[diagnose.HealthyLevel]int{
					diagnose.HealthyLevelGood:    1,
					diagnose.HealthyLevelWarn:    1,
					diagnose.HealthyLevelRisk:    1,
					diagnose.HealthyLevelSerious: 1,
					diagnose.HealthyLevelFailed:  1,
				},
			},
		},
	}); err != nil {
		t.Fatalf(err.Error())
	}
}
