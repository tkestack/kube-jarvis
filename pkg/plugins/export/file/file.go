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
package file

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

const (
	// ExporterType is type name of this Exporter
	ExporterType = "file"
)

// Exporter just print information to logger with a simple format
type Exporter struct {
	export.Collector
	*export.MetaData
	Path   string
	Format string
}

// NewExporter return a file Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	return &Exporter{
		MetaData: m,
	}
}

// Complete check and complete config items
func (e *Exporter) Complete() error {
	e.Collector.Format = e.Format
	_ = e.Collector.Complete()

	if e.Path == "" {
		e.Path = fmt.Sprintf("result.%s", e.Format)
	}

	return nil
}

// CoordinateBegin export information about coordinator Run begin
func (e *Exporter) CoordinateBegin(ctx context.Context) error {
	f, err := os.Create(e.Path)
	if err != nil {
		return errors.Wrap(err, "create file failed")
	}
	e.Collector.Output = []io.Writer{f}
	return e.Collector.CoordinateBegin(ctx)
}
