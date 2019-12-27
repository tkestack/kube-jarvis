package cron

import (
	"context"
	"testing"
	"time"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
)

func TestCoordinator_Run(t *testing.T) {
	count := 0
	c := NewCoordinator(logger.NewLogger(), fake.NewCluster()).(*Coordinator)
	f := &coordinate.FakeCoordinator{
		RunFunc: func(ctx context.Context) {
			count++
		},
	}
	c.Coordinator = f
	c.Cron = "@every 1s"

	if err := c.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	ctx, cl := context.WithTimeout(context.Background(), time.Second*10)
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
		if ctx.Err() != nil {
			t.Fatalf("timeout")
		}
		if count == 2 {
			return
		}
		time.Sleep(time.Second)
	}
}
