package translate

import "testing"

func TestTranslator_Message(t *testing.T) {
	tr, err := NewDefault("../../translation", "en", "zh")
	if err != nil {
		t.Fatalf(err.Error())
	}

	tr = tr.WithModule("diagnostics.example")
	t.Log(tr.Message("message", map[string]interface{}{
		"Mes": "test",
	}))
}
