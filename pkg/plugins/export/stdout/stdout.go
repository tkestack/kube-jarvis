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
	"io"
	"os"
	"tkestack.io/kube-jarvis/pkg/plugins/export"

	"github.com/fatih/color"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// ExporterType is type name of this Exporter
	ExporterType = "stdout"
)

// Exporter just print information to logger with a simple format
type Exporter struct {
	Format string
	export.Collector
	*export.MetaData
}

// NewExporter return a stdout Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	e := &Exporter{
		MetaData: m,
	}
	return e
}

// Complete check and complete config items
func (e *Exporter) Complete() error {
	if e.Format == "" {
		e.Format = "fmt"
	}
	e.Collector.Format = e.Format
	return e.Collector.Complete()
}

// CoordinateBegin export information about coordinator Run begin
func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	if e.Format != "fmt" {
		e.Collector.Output = []io.Writer{os.Stdout}
		return e.Collector.CoordinateBegin(ctx)
	}

	fmt.Println("===================================================================")
	fmt.Println("                       kube-jarivs                                 ")
	fmt.Println("===================================================================")
	return nil
}

// DiagnosticBegin export information about a Diagnostic begin
func (e *Exporter) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	if e.Format != "fmt" {
		return e.Collector.DiagnosticBegin(ctx, dia)
	}

	fmt.Println("Diagnostic report")
	fmt.Printf("    Type : %s\n", dia.Meta().Type)
	fmt.Printf("    Name : %s\n", dia.Meta().Name)
	fmt.Printf("- ----- results ----------------\n")
	return nil
}

// DiagnosticFinish export information about a Diagnostic finished
func (e *Exporter) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	if e.Format != "fmt" {
		return e.Collector.DiagnosticFinish(ctx, dia)
	}

	fmt.Println("===================================================================")
	return nil
}

// DiagnosticResult export information about one diagnose.Result
func (e *Exporter) DiagnosticResult(ctx context.Context, dia diagnose.Diagnostic, result *diagnose.Result) error {
	if e.Format != "fmt" {
		return e.Collector.DiagnosticResult(ctx, dia, result)
	}

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
	return nil
}

// CoordinateFinish export information about coordinator Run finished
func (e *Exporter) CoordinateFinish(ctx context.Context) error {
	if e.Format != "fmt" {
		if err := e.Collector.CoordinateFinish(ctx); err != nil {
			return err
		}
	}
	fmt.Println("===================================================================")
	return nil
}
