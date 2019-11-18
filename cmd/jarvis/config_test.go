package main

import (
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/logger"
)

func TestGetConfig(t *testing.T) {
	data := `
global:
  cluster:
    kubeconfig: "fake"
coordinate:
  type: "default"
# config:
#   xxx: "123"

diagnostics:
  - type: "example"
    name: "example 1"
    score: 10
    weight: 10
    config:
      message: "this is a example diagnostic"

evaluators:
  - type: "sum"
    name: "sum 1"
    # config:
    #   xxx: "123"

exporters:
  - type: "stdout"
    name: "stdout 1"
    # config:
    #   xxx: "123"
`

	c, err := getConfig([]byte(data), logger.NewLogger())
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
