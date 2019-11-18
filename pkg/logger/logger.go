package logger

import (
	"fmt"
	"log"
	"sort"
)

type loggerInfo struct {
	labels map[string]string
}

// NewLogger create a logger that just print logs using golang fmt logger
func NewLogger() Logger {
	return &loggerInfo{
		labels: map[string]string{},
	}
}

func (l *loggerInfo) With(labels map[string]string) Logger {
	nLogger := &loggerInfo{
		labels: map[string]string{},
	}
	for k, v := range l.labels {
		nLogger.labels[k] = v
	}

	for k, v := range labels {
		nLogger.labels[k] = v
	}

	return nLogger
}

func (l *loggerInfo) Message(prefix string, format string, args ...interface{}) string {
	message := prefix + " "
	message += fmt.Sprintf(format, args...)
	message += "  "

	keys := make([]string, 0)
	for k := range l.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		message += fmt.Sprintf("%s = %s | ", k, l.labels[k])
	}

	return message
}

func (l *loggerInfo) Infof(format string, args ...interface{}) {
	log.Println(l.Message("[INFO]", format, args...))
}
func (l *loggerInfo) Debugf(format string, args ...interface{}) {
	log.Println(l.Message("[DEBUG]", format, args...))

}
func (l *loggerInfo) Errorf(format string, args ...interface{}) {
	log.Println(l.Message("[ERROR]", format, args...))
}
