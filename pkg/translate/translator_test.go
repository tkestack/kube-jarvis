package translate

import "testing"

func TestTranslator_Message(t *testing.T) {
	_, err := NewTranslator("../../translation", "en", "zh")
	if err != nil {
		t.Fatalf(err.Error())
	}

	/*
		tr = tr.WithModule("evaluators.sum")
		t.Log(tr.Message("result", map[string]interface{}{
			"Mes": "test",
		}))*/
}
