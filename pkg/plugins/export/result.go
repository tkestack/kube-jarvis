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
package export

import (
	"encoding/json"
	"time"
	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

// DiagnosticResultItem collect one diagnostic and it's results
type DiagnosticResultItem struct {
	StartTime  time.Time
	EndTime    time.Time
	Catalogue  diagnose.Catalogue
	Type       string
	Name       string
	Desc       translate.Message
	Results    []*diagnose.Result
	Statistics map[diagnose.HealthyLevel]int
}

// NewDiagnosticResultItem return a new DiagnosticResultItem according to a Diagnostic
func NewDiagnosticResultItem(dia diagnose.Diagnostic) *DiagnosticResultItem {
	return &DiagnosticResultItem{
		StartTime:  time.Now(),
		Catalogue:  dia.Meta().Catalogue,
		Type:       dia.Meta().Type,
		Name:       dia.Meta().Name,
		Desc:       dia.Meta().Desc,
		Results:    []*diagnose.Result{},
		Statistics: map[diagnose.HealthyLevel]int{},
	}
}

func (d *DiagnosticResultItem) AddResult(r *diagnose.Result) {
	d.Results = append(d.Results, r)
	d.Statistics[r.Level]++
}

// AllResult just collect diagnostic results and progress
type AllResult struct {
	StartTime   time.Time
	EndTime     time.Time
	Diagnostics []*DiagnosticResultItem
}

// NewAllResult return a new AllResult
func NewAllResult() *AllResult {
	return &AllResult{
		StartTime:   time.Now(),
		Diagnostics: []*DiagnosticResultItem{},
	}
}

// AddDiagnosticResultItem add a diagnostic resultItem to AllResult
func (r *AllResult) AddDiagnosticResultItem(d *DiagnosticResultItem) {
	r.Diagnostics = append(r.Diagnostics, d)
}

// Marshal make AllResult become json
func (r *AllResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UnMarshal init AllResult from a json
func (r *AllResult) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
