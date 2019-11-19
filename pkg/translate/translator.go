package translate

// Message is a translated string
type Message string

// Translator translate string to target language
type Translator interface {
	// WithModule attach a module label to a Translator
	// module will be add before ID when you call Translator.Message
	Message(ID string, templateData map[string]interface{}) Message
	// Message get translated message from Translator
	// t.module will be add before ID
	// example:
	//         ID = "message"  and module = "diagnostics.example"
	//         then real ID will be "diagnostics.example.message"
	WithModule(module string) Translator
}
