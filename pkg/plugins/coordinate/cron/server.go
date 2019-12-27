package cron

import (
	"encoding/json"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"net/http"
	"tkestack.io/kube-jarvis/pkg/plugins"
)

func (c *Coordinator) runOnceHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Infof("handle run once request")
	ok := c.tryStartRun()
	if ok {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}

func (c *Coordinator) periodHandler(w http.ResponseWriter, r *http.Request) {
	// get
	if r.Method == http.MethodGet {
		if _, err := w.Write([]byte(c.Cron)); err != nil {
			c.logger.Errorf("write cron config to response failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
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

	newCron := cron.New()
	if _, err := newCron.AddFunc(string(data), c.cronDo); err != nil {
		c.logger.Errorf("create new cron failed : %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.cronLock.Lock()
	defer c.cronLock.Unlock()

	c.Cron = string(data)
	if c.cronCtl != nil {
		c.cronCtl.Stop()
	}
	c.cronCtl = newCron
	w.WriteHeader(http.StatusOK)
	c.logger.Infof("cron scheduler success update to %s", string(data))
}

type State struct {
	State    string
	Progress *plugins.Progress
}

func (c *Coordinator) stateHandler(w http.ResponseWriter, r *http.Request) {
	c.logger.Infof("handle get current state")
	resp := &State{
		Progress: c.Progress(),
	}

	if c.running {
		resp.State = "running"
	} else {
		resp.State = "pending"
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.logger.Infof("return current state success: %s ", string(data))
}
