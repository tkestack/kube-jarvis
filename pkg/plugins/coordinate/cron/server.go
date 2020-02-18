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
package cron

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"tkestack.io/kube-jarvis/pkg/httpserver"

	"github.com/robfig/cron/v3"
)

// runOnceHandler run inspection immediately
// if inspection is already running, status code will be 409
func (c *Coordinator) runOnceHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Infof("handle run once request")
	ok := c.tryStartRun()
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}

// periodHandler return or set period
// current inspection period will be returned if request method is 'Get'
// new period will be set if request method is 'POST'
func (c *Coordinator) periodHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if _, err := w.Write([]byte(c.Cron)); err != nil {
			c.logger.Errorf("write cron config to response failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	c.logger.Infof("handle update cron config")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Errorf("handle update cron config failed : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newCron := cron.New(cron.WithSeconds())
	if string(data) != "" {
		if _, err := newCron.AddFunc(string(data), c.cronDo); err != nil {
			c.logger.Errorf("create new cron failed : %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	c.cronLock.Lock()
	defer c.cronLock.Unlock()

	c.Cron = string(data)
	if c.cronCtl != nil {
		c.cronCtl.Stop()
	}

	if string(data) != "" {
		c.cronCtl = newCron
		c.cronCtl.Start()
		c.logger.Infof("cron scheduler success update to %s", string(data))
	} else {
		c.logger.Infof("cron scheduler success update closed")
	}

}

// stateHandler return current inspection state and inspection process
func (c *Coordinator) stateHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Infof("handle get current state")
	resp := &httpserver.StateResponse{
		Progress: c.Progress(),
		State:    c.state,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.logger.Errorf("marshal resp failed : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		c.logger.Errorf("write resp failed : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.logger.Infof("return current state success: %s ", string(data))
}
