package logs

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

const (
	infoLog  = "info"
	warnLog  = "warn"
	errLog   = "err"
	debug    = "debug"
	fatalLog = "fatal"
)

var (
	// PrefixFormat is the format of the prefix to every log display, includes things like time stamp and event type
	PrefixFormat = "%s %s\t"

	errColor     = color.New(color.FgRed)
	debugColor   = color.New(color.FgRed)
	warnColor    = color.New(color.FgYellow)
	defaultColor = color.New(color.FgWhite)
)

func logErrDetails(err error) {
	logPrint(debug, fmt.Sprintf("Error Details: %s", err))
}

func logPrint(logType string, message string) {
	prefix := getPrefix(logType)
	switch logType {
	case errLog:
		errColor.Println(prefix + message)
		break
	case debug:
		debugColor.Println(prefix + message)
		break
	case warnLog:
		warnColor.Println(prefix + message)
		break
	case fatalLog:
		log.Fatal(prefix + message)
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
