package rest

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
)

type Response struct {
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
	Metadata struct {
		Generated string `json:"generated"`
	} `json:"metadata"`
	Errors []ResponseError `json:"errors"`
}

type ResponseError struct {
	Severity             int    `json:"severity"`
	HumanFriendlyMessage string `json:"human_friendly_message"`
	ShortMessage         string `json:"short_message"`
}

func (r *Response) AddError(error ResponseError) {
	r.Errors = append(r.Errors, error)
}

func (r *Response) DecorateResponse(request *restful.Request) {
	r.Links.Self = "//" + request.Request.Host + request.Request.URL.String()
	r.Metadata.Generated = strconv.FormatInt(time.Now().UnixNano(), 10)
}

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
