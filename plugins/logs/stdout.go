package logs

import (
	"fmt"
	"log"

	"time"

	"github.com/fatih/color"
)

const (
	infoLog = iota
	warnLog
	errLog
	errDetailsLog
	fatalLog
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

func logPrintln(logType int, message string) {
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

func getPrefix(logType int) string {
	return fmt.Sprintf("%s %s %d\t",
		time.Now().Format("2006-01-02 15:04:05 -0700"),
		prefix,
		logType,
	)
}

func logPrintf(logType int, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	logPrintln(logType, message)
}
