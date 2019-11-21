package logger

// Logger is a abstract logger obj for kube-jarvis
type Logger interface {
	// With create a new Logger and append "labels" to old logger's labels
	With(labels map[string]string) Logger
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
