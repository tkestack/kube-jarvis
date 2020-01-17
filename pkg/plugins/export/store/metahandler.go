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
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"tkestack.io/kube-jarvis/pkg/httpserver"
)

func (e *Exporter) historyHandler(w http.ResponseWriter, r *http.Request) {
	e.hisLock.Lock()
	defer e.hisLock.Unlock()

	var err error
	var requestData []byte
	var respData []byte

	defer func() {
		e.Logger.Infof("handle meta request, err=%v, request=%s", err, string(requestData))
	}()

	defer func() { _ = r.Body.Close() }()
	requestData, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	param := &httpserver.HistoryRequest{}
	if len(requestData) != 0 {
		if err = json.Unmarshal(requestData, param); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if param.Limit == 0 {
		param.Limit = math.MaxInt32
	}

	history := httpserver.NewHistoryResponse()
	offset := param.Offset
	limit := param.Limit

	for i := len(e.history.Records) - 1; i >= 0; i-- {
		if offset != 0 {
			offset--
			continue
		}

		if limit == 0 {
			break
		}
		limit--
		history.Records = append(history.Records, e.history.Records[i])
	}

	respData, err = json.Marshal(history)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	total := 0
	n := 0
	for total < len(respData) {
		n, err = w.Write(respData[total:])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		total += n
	}
}
