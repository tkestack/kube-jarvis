package translate

type Fake struct {
}

func NewFake() Translator {
	return &Fake{}
}

// WithModule attach a module label to a Translator
// module will be add before ID when you call Translator.Message
func (f *Fake) Message(ID string, templateData map[string]interface{}) Message {
	return Message(ID)
}

// Message get translated message from Translator
// t.module will be add before ID
// example:
//         ID = "message"  and module = "diagnostics.example"
//         then real ID will be "diagnostics.example.message"
func (f *Fake) WithModule(module string) Translator {
	return f
}
