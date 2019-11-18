package main

import (
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/logger"
)

func TestGetConfig(t *testing.T) {
	c, err := GetConfig("default.yaml", logger.NewLogger())
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Logf("%+v", c)

	cli, err := c.GetClusterClient()
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = c.GetCoordinator()
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = c.GetDiagnostics(cli)
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = c.GetEvaluators()
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = c.GetExporters()
	if err != nil {
		t.Fatalf(err.Error())
	}
}
