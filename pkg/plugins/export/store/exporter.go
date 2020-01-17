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
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"tkestack.io/kube-jarvis/pkg/httpserver"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

const (
	// ExporterType is type name of this Exporter
	ExporterType     = "store"
	resultsStoreName = "results"
)

// Exporter save result to file
type Exporter struct {
	*export.MetaData
	// MaxRemain
	MaxRemain int
	Path      string
	Server    bool

	history *export.History
	hisLock sync.Mutex
}

// NewExporter return a file Exporter
func NewExporter(m *export.MetaData) export.Exporter {
	e := &Exporter{
		MetaData: m,
		history:  &export.History{},
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
		httpserver.Default.HandleFunc(httpserver.StandardQueryPath, e.queryHandler)
		httpserver.Default.HandleFunc(httpserver.StandardHistoryPath, e.historyHandler)
	}

	if _, err := e.Store.CreateSpace(resultsStoreName); err != nil {
		return err
	}
	return e.reloadHistory()
}

func (e *Exporter) reloadHistory() error {
	e.hisLock.Lock()
	defer e.hisLock.Unlock()
	data, _, err := e.Store.Get(resultsStoreName, "history")
	if err != nil {
		return nil
	}

	if data != "" {
		return json.Unmarshal([]byte(data), e.history)
	}
	return nil
}

// Export export result
func (e *Exporter) Export(ctx context.Context, result *export.AllResult) error {
	e.hisLock.Lock()
	defer e.hisLock.Unlock()

	ID := fmt.Sprint(result.StartTime.UnixNano())
	if err := e.saveResult(ID, result); err != nil {
		return err
	}

	// create meta item
	e.history.Records = append(e.history.Records, &export.HistoryItem{
		ID: ID,
		Overview: export.AllResult{
			StartTime:  result.StartTime,
			EndTime:    result.EndTime,
			Statistics: result.Statistics,
		},
	})

	if e.MaxRemain >= len(e.history.Records) {
		return e.saveHistory()
	}

	// remove old items
	index := len(e.history.Records) - e.MaxRemain
	removeItems := e.history.Records[0:index]
	e.history.Records = e.history.Records[index:]
	if err := e.saveHistory(); err != nil {
		return err
	}

	for _, item := range removeItems {
		_ = e.removeResult(item.ID)
	}

	return nil
}

func (e *Exporter) saveResult(ID string, result *export.AllResult) error {
	data, err := result.Marshal()
	if err != nil {
		return err
	}

	return e.Store.Set(resultsStoreName, ID, string(data))
}

func (e *Exporter) saveHistory() error {
	data, err := e.history.Marshal()
	if err != nil {
		return err
	}

	return e.Store.Set(resultsStoreName, "history", string(data))
}

func (e *Exporter) removeResult(ID string) error {
	return e.Store.Delete(resultsStoreName, ID)
}
