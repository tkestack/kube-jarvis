package file

import (
	"context"
	"fmt"
	"os"
	"sync"

	"tkestack.io/kube-jarvis/pkg/httpserver"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

const (
	// ExporterType is type name of this Exporter
	ExporterType = "file"
)

// Exporter save result to file
type Exporter struct {
	MaxRemain int
	Path      string
	Server    bool
	*export.MetaData
	meta     *Meta
	metaLock sync.Mutex
}

// NewExporter return a file Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	e := &Exporter{
		MetaData: m,
	}
	return e
}

// Complete check and complete config items
func (e *Exporter) Complete() error {
	if e.Path == "" {
		e.Path = "results"
	}

	if e.MaxRemain == 0 {
		e.MaxRemain = 7
	}

	if e.Server {
		httpserver.HandleFunc("/exporter/file/query", e.queryHandler)
		httpserver.HandleFunc("/exporter/file/meta", e.metaHandler)
	}
	return e.reloadMeta()
}

// Export export result
func (e *Exporter) Export(ctx context.Context, result *export.AllResult) error {
	e.metaLock.Lock()
	defer e.metaLock.Unlock()

	_ = os.MkdirAll(e.Path, 0755)
	if err := e.reloadMeta(); err != nil {
		return err
	}

	ID := fmt.Sprint(result.StartTime.UnixNano())

	// save file
	if err := e.saveResult(ID, result); err != nil {
		return err
	}

	// create meta item
	e.meta.Results = append(e.meta.Results, &MetaItem{
		ID: ID,
		Overview: export.AllResult{
			StartTime:  result.StartTime,
			EndTime:    result.EndTime,
			Statistics: result.Statistics,
		},
	})

	if e.MaxRemain >= len(e.meta.Results) {
		return e.saveMeta()
	}

	// remove old items
	index := len(e.meta.Results) - e.MaxRemain
	removeItems := e.meta.Results[0:index]
	e.meta.Results = e.meta.Results[index:]
	if err := e.saveMeta(); err != nil {
		return err
	}

	for _, item := range removeItems {
		fileName := fmt.Sprintf("%s/%s.json", e.Path, item.ID)
		if err := os.Remove(fileName); err != nil {
			e.Logger.Errorf("remove data file %s failed:%v", fileName, err)
		}
	}

	return nil
}
