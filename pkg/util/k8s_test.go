package util

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
)

func Test_Quantity(t *testing.T) {
	a := resource.MustParse("7.63Gi")
	if MemQuantityStr(&a) != "7.63GB" {
		t.Fatalf("shoud be 7.63GB")
	}

	a = resource.MustParse("1100m")
	t.Log(CpuQuantityStr(&a))
}
