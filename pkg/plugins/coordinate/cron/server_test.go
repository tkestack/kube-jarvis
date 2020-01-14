package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"tkestack.io/kube-jarvis/pkg/httpserver"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster/fake"
	"tkestack.io/kube-jarvis/pkg/plugins/coordinate"
	"tkestack.io/kube-jarvis/pkg/store"
)

func Test_runOnceHandler(t *testing.T) {
	var cases = []struct {
		running    bool
		returnCode int
	}{
		{
			running:    false,
			returnCode: http.StatusOK,
		},
		{
			running:    true,
			returnCode: http.StatusConflict,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			c := NewCoordinator(logger.NewLogger(), fake.NewCluster(), store.GetStore("mem")).(*Coordinator)
			ctx, cl := context.WithCancel(context.Background())
			defer cl()
			c.Coordinator = &coordinate.FakeCoordinator{RunFunc: func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			}}
			go func() {
				_ = c.Run(ctx)
			}()

			if cs.running {
				c.tryStartRun()
			}

			resp := httpserver.NewFakeResponseWriter()
			c.runOnceHandler(resp, nil)
			if resp.StatusCode != cs.returnCode {
				t.Fatalf("statusCode want %d but get %d", cs.returnCode, resp.StatusCode)
			}
		})
	}
}

type progressCoordinator struct {
	coordinate.FakeCoordinator
	progress *plugins.Progress
}

func (p *progressCoordinator) Progress() *plugins.Progress {
	return p.progress
}

func Test_state(t *testing.T) {
	c := NewCoordinator(logger.NewLogger(), fake.NewCluster(), store.GetStore("mem")).(*Coordinator)
	c.state = StateRunning
	co := &progressCoordinator{
		progress: plugins.NewProgress(),
	}
	co.progress.Percent = 100
	c.Coordinator = co

	req := &http.Request{
		Method: http.MethodGet,
	}

	resp := httpserver.NewFakeResponseWriter()
	c.stateHandler(resp, req)

	result := &State{}
	if err := json.Unmarshal(resp.RespData, result); err != nil {
		t.Fatalf(err.Error())
	}

	if c.state != result.State {
		t.Fatalf("want running")
	}

	if result.Progress.Percent != co.progress.Percent {
		t.Fatalf("want percent 100")
	}
}

func Test_getCron(t *testing.T) {
	c := NewCoordinator(logger.NewLogger(), fake.NewCluster(), store.GetStore("mem")).(*Coordinator)
	c.Cron = "1 1 1 1 1"
	req := &http.Request{
		Method: http.MethodGet,
	}
	resp := httpserver.NewFakeResponseWriter()
	c.periodHandler(resp, req)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("should return statusOK but get %d", resp.StatusCode)
	}

	if string(resp.RespData) != c.Cron {
		t.Fatalf("want %s but get %s", c.Cron, string(resp.RespData))
	}
}

func Test_updateCron(t *testing.T) {
	var cases = []struct {
		cron       string
		statusCode int
	}{
		{
			cron:       "1 1 1 1 1 1",
			statusCode: http.StatusOK,
		},
		{
			cron:       "xx",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			c := NewCoordinator(logger.NewLogger(), fake.NewCluster(), store.GetStore("mem")).(*Coordinator)
			resp := httpserver.NewFakeResponseWriter()
			req := &http.Request{
				Method: http.MethodPut,
				Body:   ioutil.NopCloser(strings.NewReader(cs.cron)),
			}
			c.periodHandler(resp, req)
			if resp.StatusCode != cs.statusCode {
				t.Fatalf("want status %d but get %d", cs.statusCode, resp.StatusCode)
			}

			if resp.StatusCode != http.StatusOK {
				return
			}

			if cs.cron != c.Cron {
				t.Fatalf("update Cron failed")
			}
		})
	}
}
