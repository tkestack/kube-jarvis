package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

type QueryRequest struct {
	ID     string
	Type   string
	Name   string
	Level  diagnose.HealthyLevel
	Offset int
	Limit  int
}

func (e *Exporter) queryHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var requestData []byte
	var respData []byte

	defer func() {
		e.Logger.Infof("handle query request, err=%v, request=%s", err, string(requestData))
	}()

	defer func() { _ = r.Body.Close() }()
	requestData, err = ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	param := &QueryRequest{}
	if err = json.Unmarshal(requestData, param); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if param.Level != "" && !param.Level.Verify() {
		err = fmt.Errorf("unknown 'Level'='%s'", param.Level)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if param.Limit == 0 {
		param.Limit = math.MaxInt32
	}

	content, exist, err := e.Store.Get(resultsStoreName, param.ID)
	if err != nil {
		err = fmt.Errorf("get result failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !exist {
		err = fmt.Errorf("result not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	allResults := export.NewAllResult()
	if err = allResults.UnMarshal([]byte(content)); err != nil {
		return
	}

	respResult := export.NewAllResult()
	respResult.StartTime = allResults.StartTime
	respResult.EndTime = allResults.EndTime
	respResult.Statistics = allResults.Statistics

	for _, dia := range allResults.Diagnostics {
		if param.Type != "" && param.Type != dia.Type {
			continue
		}

		if param.Name != "" && param.Name != dia.Name {
			continue
		}

		newDia := &export.DiagnosticResultItem{
			StartTime:  dia.StartTime,
			EndTime:    dia.EndTime,
			Catalogue:  dia.Catalogue,
			Type:       dia.Type,
			Name:       dia.Name,
			Desc:       dia.Desc,
			Results:    []*diagnose.Result{},
			Statistics: dia.Statistics,
		}

		offset := param.Offset
		limit := param.Limit
		for _, item := range dia.Results {
			if param.Level != "" && item.Level.Compare(param.Level) > 0 {
				continue
			}

			if offset > 0 {
				offset--
				continue
			}

			if limit == 0 {
				break
			}

			newDia.Results = append(newDia.Results, item)
		}
		respResult.Diagnostics = append(respResult.Diagnostics, newDia)
	}

	respData, err = json.Marshal(respResult)
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
