package translate

import "testing"

func TestNewFake(t *testing.T) {
	f := NewFake()
	f.WithModule("123")
	ms := f.Message("123", nil)
	if ms != "123" {
		t.Fatalf("return msg should be 123")
	}
}
