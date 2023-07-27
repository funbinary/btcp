package logger

import "fmt"

type Logger struct {
}

func (l *Logger) Trace(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Tracef(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Logger) Debug(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Debugf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Logger) Info(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Infof(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (l *Logger) Warn(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Warnf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (l *Logger) Error(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Errorf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
func (l *Logger) Fatal(args ...interface{}) {
	fmt.Println(args...)
}
func (l *Logger) Fatalf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
