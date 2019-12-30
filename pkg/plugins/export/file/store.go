package file

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"tkestack.io/kube-jarvis/pkg/plugins/export"
)

type MetaItem struct {
	ID       string
	Overview export.AllResult
}

type Meta struct {
	Results []*MetaItem
}

func (e *Exporter) reloadMeta() error {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", e.Path, "meta.json"))
	if err != nil {
		if os.IsNotExist(err) {
			e.meta = &Meta{Results: []*MetaItem{}}
			return nil
		}
		return errors.Wrapf(err, "read file failed")
	}

	meta := &Meta{Results: []*MetaItem{}}
	if err := json.Unmarshal(data, meta); err != nil {
		return errors.Wrapf(err, "json unmarshal failed")
	}

	e.meta = meta
	return nil
}

func (e *Exporter) saveMeta() error {
	data, err := json.Marshal(e.meta)
	if err != nil {
		return errors.Wrapf(err, "marshal meta data failed")
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/%s", e.Path, "meta.json"), data, 0666)
}

func (e *Exporter) saveResult(ID string, result *export.AllResult) error {
	data, err := result.Marshal()
	if err != nil {
		return errors.Wrapf(err, "marshal result failed")
	}

	fileName := fmt.Sprintf("%s/%s.json", e.Path, ID)
	if err := ioutil.WriteFile(fileName, data, 0666); err != nil {
		return errors.Wrapf(err, "write file failed")
	}
	return nil
}
