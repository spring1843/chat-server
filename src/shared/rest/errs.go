package rest

import (
	"fmt"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
)

// Syslog model error severity levels https://en.wikipedia.org/wiki/Syslog#Severity_level
const (
	// Alert should be corrected immediately e.g. Loss of the primary ISP connection.
	Alert = "alert"

	// Critical conditions e.g. A failure in the system's primary application.
	Critical = "crit"

	// Error conditions e.g. An application has exceeded its file storage limit and attempts to write are failing.
	Error = "err"

	// Warning may indicate that an error will occur if action is not taken. e.g. A non-root file system has only 2GB remaining.
	Warning = "warn"

	// Notice events that are unusual, but not error conditions. e.g. Ski Haus Delta reports temperature < low_notice(50)
	Notice = "notice"

	// Informational normal operational messages that require no action. e.g. An application has started, paused or ended successfully.
	Informational = "info"

	// Debug information useful to developers for debugging the application.
	Debug = "debug"
)

// RespError is an error that is part of the API response
type RespError struct {
	Severity             string `json:"severity"`
	HumanFriendlyMessage string `json:"human_friendly_message"`
	ShortMessage         string `json:"short_message"`
	Details              string `json:"details"`
}

// AddError adds a REST error to response
func (r *Resp) AddError(restError RespError) {
	r.Errors = append(r.Errors, restError)
}

// AddDetailedError adds a formatted REST error to response
func (r *Resp) AddDetailedError(restError RespError, format string, a ...interface{}) {
	restError.Details = fmt.Sprintf(format, a...)
	r.Errors = append(r.Errors, restError)
}

// DecorateResponse is executed before a response is generated
func (r *Resp) DecorateResponse(request *restful.Request) {
	r.Links.Self = "//" + request.Request.Host + request.Request.URL.String()
	r.Metadata.Generated = strconv.FormatInt(time.Now().UnixNano(), 10)
}
