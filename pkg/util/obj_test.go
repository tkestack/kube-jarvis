package util

import (
	"testing"
)

func TestInitObjViaYaml(t *testing.T) {
	type E struct {
		A string
	}
	type T struct {
		E
	}

	var a T
	var b T
	a.A = "123"

	if err := InitObjViaYaml(&b, &a); err != nil {
		t.Fatalf(err.Error())
	}

	if b.A != a.A {
		t.Fatalf("b.A != a.A")
	}
}
