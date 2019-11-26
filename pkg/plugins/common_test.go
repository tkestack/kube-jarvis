package plugins

import (
	"fmt"
	"testing"
)

func TestIsSupportedCloud(t *testing.T) {
	var cases = []struct {
		supported bool
		clouds    []string
		cloud     string
	}{
		{
			supported: true,
			clouds:    []string{},
			cloud:     "123",
		},
		{
			supported: true,
			clouds: []string{
				"123",
			},
			cloud: "123",
		},
		{
			supported: false,
			clouds: []string{
				"321",
			},
			cloud: "123",
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			if IsSupportedCloud(cs.clouds, cs.cloud) != cs.supported {
				t.Fatalf("shoud %v", cs.supported)
			}
		})
	}
}
