package store

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"

	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

type HistoryRequest struct {
	Offset int
	Limit  int
}

func (e *Exporter) metaHandler(w http.ResponseWriter, r *http.Request) {
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

	param := &HistoryRequest{}
	if len(requestData) != 0 {
		if err = json.Unmarshal(requestData, param); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if param.Limit == 0 {
		param.Limit = math.MaxInt32
	}

	history := &export.History{
		Records: []*export.HistoryItem{},
	}

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
