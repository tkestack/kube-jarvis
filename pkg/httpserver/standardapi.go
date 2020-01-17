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
package httpserver

import (
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

const (
	// StandardQueryPath is the standard API path for querying target diagnostic results
	StandardQueryPath = "/exporter/store/query"
	// StandardHistoryPath is the standard API path for querying diagnostic history
	StandardHistoryPath = "/exporter/store/history"
	// StandardRunPath is the standard API path for starting diagnose immediately
	StandardRunPath = "/coordinator/cron/run"
	// StandardStatePath is the standard API path for getting current running state and progress
	// this API only available when the coordinator type is "cron"
	StandardStatePath = "/coordinator/cron/state"
	// StandardPeriodPath is the standard path for getting or updating running period
	// this API only available when the coordinator type is "cron"
	StandardPeriodPath = "/coordinator/cron/period"
)

// HistoryRequest is the request for querying history records
type HistoryRequest struct {
	// Offset is the offset of target records
	// the histories will be sorted in descending chronological order
	// so, offset=1 means the second recent record
	Offset int
	// Limit is the max record number of returned records
	Limit int
}

// HistoryResponse is the response of query history records
type HistoryResponse struct {
	*export.History
}

// NewHistoryResponse return a HistoryResponse with default values
func NewHistoryResponse() *HistoryResponse {
	return &HistoryResponse{
		&export.History{Records: []*export.HistoryItem{}},
	}
}

// QueryRequest is the request for querying one diagnostic report
type QueryRequest struct {
	// ID used to specify the target report
	ID string
	// Type is the target type of diagnostic
	// if Type is not empty, only diagnostic with type "Type" will be returned
	Type string
	// Name is the target name of diagnostic
	// if Name is not empty, only diagnostic with name "Name" will be returned
	Name string
	// Level is the max HealthyLevel of target results
	// for example:
	//      Level = HealthyLevelSerious
	//      will only return results with HealthyLevel HealthyLevelSerious, and HealthyLevelFailed
	// if Level is empty ,HealthyLevelGood will be used
	Level diagnose.HealthyLevel
	// Offset is the offset value of the request result
	Offset int
	// Limit is the max line of results
	Limit int
}

// QueryResponse is the response of querying results
type QueryResponse struct {
	*export.AllResult
}

// NewQueryResponse create an empty QueryResponse
func NewQueryResponse() *QueryResponse {
	return &QueryResponse{
		AllResult: export.NewAllResult(),
	}
}

// StateResponse is the response of querying current state
type StateResponse struct {
	State    string
	Progress *plugins.Progress
}
