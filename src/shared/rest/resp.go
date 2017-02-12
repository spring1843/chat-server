package rest

import (
	"encoding/json"
	"io"

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
)

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
