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
package translate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// Default translate string to target language
type Default struct {
	module   string
	bundle   *i18n.Bundle
	localize *i18n.Localizer
}

// NewDefault create a new default Translator
// Translator will read translation message from "dir/defLang" and "dir/targetLang"
func NewDefault(dir string, defLang string, targetLang string) (Translator, error) {
	t := &Default{}
	defTag, err := language.Parse(defLang)
	if err != nil {
		return nil, err
	}

	targetTag := language.Make(targetLang)
	t.bundle = i18n.NewBundle(defTag)
	t.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	t.localize = i18n.NewLocalizer(t.bundle, targetLang)

	// load default message
	if err := t.addMessage(dir, defTag); err != nil {
		return nil, err
	}

	// load target message
	return t, t.addMessage(dir, targetTag)
}

func (d *Default) addMessage(dir string, tag language.Tag) error {
	return filepath.Walk(fmt.Sprintf("%s/%s", dir, tag.String()), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".yaml") {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			mes, err := i18n.ParseMessageFileBytes(buf, path, map[string]i18n.UnmarshalFunc{
				"yaml": yaml.Unmarshal,
			})
			if err != nil {
				return fmt.Errorf("load message file %s failed : %s", path, err.Error())
			}

			for _, m := range mes.Messages {
				m.ID = fmt.Sprintf("%s.%s.%s", filepath.Base(filepath.Dir(path)), strings.TrimSuffix(info.Name(), ".yaml"), m.ID)
			}

			if err := d.bundle.AddMessages(tag, mes.Messages...); err != nil {
				return fmt.Errorf("add message failed : %s", err.Error())
			}
		}
		return nil
	})
}

// WithModule attach a module label to a Translator
// module will be add before ID when you call Translator.Message
func (d *Default) WithModule(module string) Translator {
	return &Default{
		module:   module,
		bundle:   d.bundle,
		localize: d.localize,
	}
}

// Message get translated message from Translator
// t.module will be add before ID
// example:
//         ID = "message"  and module = "diagnostics.example"
//         then real ID will be "diagnostics.example.message"
func (d *Default) Message(ID string, templateData map[string]interface{}) Message {
	mes, _ := d.localize.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: fmt.Sprintf("%s.%s", d.module, ID),
		},
		TemplateData: templateData,
	})
	return Message(mes)
}
