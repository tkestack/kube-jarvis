package main

import (
	"context"
	"flag"

	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	_ "github.com/RayHuangCN/kube-jarvis/pkg/plugins/coordinate/all"
	_ "github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose/all"
	_ "github.com/RayHuangCN/kube-jarvis/pkg/plugins/evaluate/all"
	_ "github.com/RayHuangCN/kube-jarvis/pkg/plugins/export/all"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "conf/default.yaml", "config file")
}

func main() {
	config, err := GetConfig(configFile, logger.NewLogger())
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

	trans, err := config.GetTranslator()
	if err != nil {
		panic(err)
	}

	diagnostics, err := config.GetDiagnostics(cli, trans)
	if err != nil {
		panic(err)
	}

	for _, d := range diagnostics {
		coordinator.AddDiagnostic(d)
	}

	evaluators, err := config.GetEvaluators(cli, trans)
	if err != nil {
		panic(err)
	}

	for _, e := range evaluators {
		coordinator.AddEvaluate(e)
	}

	exporters, err := config.GetExporters(cli)
	if err != nil {
		panic(err)
	}

	for _, e := range exporters {
		coordinator.AddExporter(e)
	}

	coordinator.Run(context.Background())
}
