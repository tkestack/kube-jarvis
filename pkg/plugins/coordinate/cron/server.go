package cron

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/robfig/cron/v3"
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
	c.cronCtl.Start()
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
