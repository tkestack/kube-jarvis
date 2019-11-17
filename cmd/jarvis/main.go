package main

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/logger"

	"github.com/RayHuangCN/Jarvis/cmd/jarvis/conf"
)

func main() {
	config, err := conf.GetConfig("conf/default.yaml", logger.NewLogger())
	if err != nil {
		panic(err)
	}

	cli, err := config.GetClusterClient()
	if err != nil {
		panic(err)
	}

	coordinator, err := config.GetCoordinator()
	if err != nil {
		panic(err)
	}

	diagnostics, err := config.GetDiagnostics(cli)
	if err != nil {
		panic(err)
	}

	for _, d := range diagnostics {
		coordinator.AddDiagnostic(d)
	}

	evaluators, err := config.GetEvaluators()
	if err != nil {
		panic(err)
	}

	for _, e := range evaluators {
		coordinator.AddEvaluate(e)
	}

	exporters, err := config.GetExporters()
	if err != nil {
		panic(err)
	}

	for _, e := range exporters {
		coordinator.AddExporter(e)
	}

	coordinator.Run(context.Background())
}
