package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	lg := NewLogger()
	lg = lg.With(map[string]string{
		"user": "test",
	}).With(map[string]string{
		"logger": "fmt",
	})

	lg.Infof("info")
	lg.Errorf("error")
	lg.Debugf("debug")
}
