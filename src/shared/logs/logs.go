package logs

import "fmt"

// Simple message

// Info info logs
func Info(message string) {
	logPrint(infoLog, message)
}

// Fatal fatal logs
func Fatal(message string) {
	logPrint(fatalLog, message)
}

// Warn warn logs
func Warn(message string) {
	logPrint(warnLog, message)
}

// Err error logs
func Err(message string) {
	logPrint(errLog, message)
}

// With format

// Infof info logs if format
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a))
}

// Fatalf fatal logs with format
func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a))
}

// Warnf warn logs with format
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a))
}

// Errf error logs with format
func Errf(format string, a ...interface{}) {
	Err(fmt.Sprintf(format, a))
}

// With error

// FatalIfErrf logs fatal if there's an error
func FatalIfErrf(err error, format string, a ...interface{}) {
	if err != nil {
		logErrDetails(err)
		Fatalf(format, a...)
	}
}

// WarnIfErrf warn logs if there's an error
func WarnIfErrf(err error, format string, a ...interface{}) {
	if err != nil {
		logErrDetails(err)
		Warnf(format, a...)
	}
}

// ErrIfErrf error logs if there's an error
func ErrIfErrf(err error, format string, a ...interface{}) {
	if err != nil {
		logErrDetails(err)
		Errf(format, a...)
	}
}
