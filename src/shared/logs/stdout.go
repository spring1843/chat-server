package logs

import (
	"fmt"
	"log"
	"time"

	"github.com/spring1843/chat-server/libs/color"
	"github.com/spring1843/chat-server/libs/errors"
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
	PrefixFormat = "%s %s "

	infoColor    = color.New(color.FgHiWhite)
	errColor     = color.New(color.FgRed)
	debugColor   = color.New(color.FgCyan)
	warnColor    = color.New(color.FgYellow)
	defaultColor = color.New(color.FgWhite)
)

func logErrDetails(err error) {
	type causer interface {
		Cause() error
	}
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
		logPrint(debug, fmt.Sprintf("Caused By: %s", err))
		logPrint(debug, fmt.Sprintf("Stack: %+v", errors.WithStack(err)))
	}
}

func logPrint(logType string, message string) {
	prefix := getPrefix(logType)
	switch logType {
	case infoLog:
		infoColor.Println(prefix + message)
		break
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
