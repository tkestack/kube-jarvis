package file

import (
	"context"
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
	Path string
}

// NewExporter return a file Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	return &Exporter{
		MetaData: m,
	}
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
