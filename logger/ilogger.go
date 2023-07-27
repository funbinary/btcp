package logger

type ILogger interface {
	Trace(args ...interface{})
	Tracef(format string, v ...interface{})
	Debug(args ...interface{})
	Debugf(format string, v ...interface{})
	Info(args ...interface{})
	Infof(format string, v ...interface{})
	Warn(args ...interface{})
	Warnf(format string, v ...interface{})
	Error(args ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, v ...interface{})
}
