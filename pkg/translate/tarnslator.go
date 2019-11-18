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

// Message is a translated string
type Message string

// Translator translate string to target language
type Translator struct {
	module   string
	bundle   *i18n.Bundle
	localize *i18n.Localizer
}

// NewTranslator create a new Translator
// Translator will read translation message from "dir/defLang" and "dir/targetLang"
func NewTranslator(dir string, defLang string, targetLang string) (*Translator, error) {
	t := &Translator{}
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

func (t *Translator) addMessage(dir string, tag language.Tag) error {
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

			if err := t.bundle.AddMessages(tag, mes.Messages...); err != nil {
				return fmt.Errorf("add message failed : %s", err.Error())
			}
		}
		return nil
	})
}

// WithModule attach a module label to a Translator
// module will be add before ID when you call Translator.Message
func (t *Translator) WithModule(module string) *Translator {
	return &Translator{
		module:   module,
		bundle:   t.bundle,
		localize: t.localize,
	}
}

// Message get translated message from Translator
// t.module will be add before ID
// example:
//         ID = "message"  and module = "diagnostics.example"
//         then real ID will be "diagnostics.example.message"
func (t *Translator) Message(ID string, templateData map[string]interface{}) Message {
	return Message(t.localize.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: fmt.Sprintf("%s.%s", t.module, ID),
		},
		TemplateData: templateData,
	}))
}
