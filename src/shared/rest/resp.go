package rest

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/spring1843/chat-server/libs/go-restful"
)

type (
	// Resp is the structure that will be in every API response
	Resp struct {
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Metadata struct {
			Generated string `json:"generated"`
		} `json:"metadata"`
		Errors []RespError `json:"errors"`
	}
	// RespError is an error that is part of the API response
	RespError struct {
		Severity             int    `json:"severity"`
		HumanFriendlyMessage string `json:"human_friendly_message"`
		ShortMessage         string `json:"short_message"`
	}
)

// AddError adds a REST error to response
func (r *Resp) AddError(restError RespError) {
	r.Errors = append(r.Errors, restError)
}

// DecorateResponse is executed before a response is generated
func (r *Resp) DecorateResponse(request *restful.Request) {
	r.Links.Self = "//" + request.Request.Host + request.Request.URL.String()
	r.Metadata.Generated = strconv.FormatInt(time.Now().UnixNano(), 10)
}

// ParseRequestBody parses request  body
func ParseRequestBody(r *restful.Request, o interface{}) error {
	raw := make([]byte, r.Request.ContentLength)
	_, err := io.ReadFull(r.Request.Body, raw)
	if err != nil {
		return err
	}
	err = json.Unmarshal(raw, o)
	if err != nil {
		return err
	}
	return nil
}
