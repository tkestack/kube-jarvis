package cron

import (
	"context"
	"testing"
	"time"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

type fakeCoordinator struct {
	RunCount int
}

// Complete check and complete config items
func (f *fakeCoordinator) Complete() error {
	return nil
}

// AddDiagnostic add a diagnostic to Coordinator
func (f *fakeCoordinator) AddDiagnostic(dia diagnose.Diagnostic) {

}

// AddExporter add a Exporter to Coordinator
func (f *fakeCoordinator) AddExporter(exporter export.Exporter) {

}

// Run will do all diagnostics, evaluations, then export it by exporters
func (f *fakeCoordinator) Run(ctx context.Context) {
	f.RunCount++
}

func (f *fakeCoordinator) Progress() *plugins.Progress {
	return nil
}

func TestCoordinator_Run(t *testing.T) {
	c := NewCoordinator(logger.NewLogger(), fake.NewCluster()).(*Coordinator)
	f := &fakeCoordinator{}
	c.Coordinator = f

	ctx, cl := context.WithCancel(context.Background())
	defer cl()
	go func() {
		c.Run(ctx)
	}()

	for {
		suc := c.tryStartRun()
		if suc {
			break
		}
	}
	for {
		suc := c.tryStartRun()
		if suc {
			break
		}
	}
	time.Sleep(time.Second)
	if f.RunCount != 2 {
		t.Fatalf("should run 2")
	}

}
