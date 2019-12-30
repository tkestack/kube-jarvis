package file

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
)

type MetaRequest struct {
	Offset int
	Limit  int
}

func (e *Exporter) metaHandler(w http.ResponseWriter, r *http.Request) {
	e.metaLock.Lock()
	defer e.metaLock.Unlock()

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

	param := &MetaRequest{}
	if len(requestData) != 0 {
		if err = json.Unmarshal(requestData, param); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if param.Limit == 0 {
		param.Limit = math.MaxInt32
	}

	meta := &Meta{
		Results: []*MetaItem{},
	}

	offset := param.Offset
	limit := param.Limit

	for i := len(e.meta.Results) - 1; i >= 0; i-- {
		if offset != 0 {
			offset--
			continue
		}

		if limit == 0 {
			break
		}
		limit--
		meta.Results = append(meta.Results, e.meta.Results[i])
	}

	respData, err = json.Marshal(meta)
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
