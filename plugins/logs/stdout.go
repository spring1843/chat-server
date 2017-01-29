package logs

import (
	"fmt"
	"log"

	"time"

	"github.com/fatih/color"
)

const (
	infoLog       = "info"
	warnLog       = "warn"
	errLog        = "err"
	errDetailsLog = "err_details"
	fatalLog      = "fatal"
)

var (
	errColor        = color.New(color.FgRed)
	errDetailsColor = color.New(color.FgRed)
	warnColor       = color.New(color.FgYellow)
	defaultColor    = color.New(color.FgWhite)
)

func logFatal(err error) {
	log.Fatalf("Fatal Error: %s", err)
}

func logErrDetails(err error) {
	logPrintf(errDetailsLog, "Error Details: %s", err.Error())
}

func logPrintf(logType string, format string, a ...interface{}) {
	logPrint(logType, fmt.Sprintf(format, a...))
}

func logPrint(logType string, message string) {
	logPrintln(logType, message)
}

func logPrintln(logType string, message string) {
	prefix := getPrefix(logType)
	switch logType {
	case errLog:
		errColor.Println(prefix + message)
		break
	case errDetailsLog:
		errDetailsColor.Println(prefix + message)
		break
	case warnLog:
		warnColor.Println(prefix + message)
		break
	default:
		defaultColor.Println(prefix + message)
	}
}

func getPrefix(logType string) string {
	return fmt.Sprintf(
		PrefixFormat,
		time.Now().Format("2006-01-02 15:04:05 -0700"),
		logType,
	)
}
