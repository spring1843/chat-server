package rest

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spring1843/chat-server/libs/go-restful"
)

// RespError is an error that is part of the API response
type RespError struct {
	Severity             int    `json:"severity"`
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
