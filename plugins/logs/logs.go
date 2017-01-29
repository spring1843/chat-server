package logs

var PrefixFormat = "%s %s\t"

// Infof Logs formatted information
func Infof(format string, a ...interface{}) {
	logPrintf(infoLog, format, a...)
}

// Infof Logs information
func Info(message string) {
	logPrint(infoLog, message)
}

// Infof Logs formatted information
func Fatalf(format string, a ...interface{}) {
	logPrintf(fatalLog, format, a...)
}

// Warnf logs warnings, and an error
func Warnf(err error, format string, a ...interface{}) {
	if err != nil {
		logPrintf(warnLog, format, a...)
		logErrDetails(err)
	}
}

// Errf logs errir and an error
func Errf(err error, format string, a ...interface{}) {
	if err != nil {
		logPrintf(errLog, format, a...)
		logErrDetails(err)
	}
}

// FatalErrf logs a fatal message, and error and ends execution
func FatalErrf(err error, format string, a ...interface{}) {
	if err != nil {
		logPrintf(fatalLog, format, a...)
		logFatal(err)
	}
}
