/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package stdout

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"tkestack.io/kube-jarvis/pkg/plugins/export"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// ExporterType is type name of this Exporter
	ExporterType = "stdout"
)

// Exporter just print information to logger with a simple format
type Exporter struct {
	Format string
	Level  diagnose.HealthyLevel
	*export.Collector
	*export.MetaData
}

// NewExporter return a stdout Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	e := &Exporter{
		MetaData:  m,
		Collector: export.NewCollector(),
	}
	return e
}

// Complete check and complete config items
func (e *Exporter) Complete() error {
	if e.Format == "" {
		e.Format = "fmt"
	}
	e.Collector.Format = e.Format

	if e.Level == "" {
		e.Level = diagnose.HealthyLevelGood
	}

	if !e.Level.Verify() {
		return fmt.Errorf("level %s is illegal", e.Level)
	}
	e.Collector.Level = e.Level

	return e.Collector.Complete()
}

// CoordinateBegin export information about coordinator Run begin
func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	if e.Format != "fmt" {
		e.Collector.Output = []io.Writer{os.Stdout}
	}

	return e.Collector.CoordinateBegin(ctx)
}

// CoordinateFinish export information about coordinator Run finished
func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	if e.Format != "fmt" {
		if err := e.Collector.CoordinateFinish(ctx); err != nil {
			return err
		}
	}

	fmt.Println("===================================================================")
	fmt.Println("                       kube-jarivs                                 ")
	fmt.Println("===================================================================")

	for _, dia := range e.Collector.Diagnostics {
		fmt.Println("Diagnostic report")
		fmt.Printf("    Type : %s\n", dia.Type)
		fmt.Printf("    Desc : %s\n", dia.Desc)
		fmt.Printf("    Name : %s\n", dia.Name)
		fmt.Printf("    TotalResult  : %d\n", dia.Total)
		fmt.Printf("    CollectedResult  : %d\n", len(dia.Results))
		fmt.Printf("- ----- results ----------------\n")

		for _, result := range dia.Results {
			var pt func(format string, a ...interface{})
			switch result.Level {
			case diagnose.HealthyLevelFailed:
				pt = color.HiRed
			case diagnose.HealthyLevelGood:
				pt = color.Green
			case diagnose.HealthyLevelWarn:
				pt = color.Yellow
			case diagnose.HealthyLevelRisk:
				pt = color.Red
			case diagnose.HealthyLevelSerious:
				pt = color.HiRed
			default:
				pt = func(format string, a ...interface{}) {
					fmt.Printf(format, a...)
				}
			}
			pt("[%s] %s -> %s\n", result.Level, result.Title, result.ObjName)
			pt("    Describe : %s\n", result.Desc)
			pt("    Proposal : %s\n", result.Proposal)
			fmt.Printf("- -----------------------------\n")
		}
	}

	fmt.Println("===================================================================")
	return nil
}
