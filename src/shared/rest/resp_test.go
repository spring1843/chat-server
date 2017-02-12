package rest_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/spring1843/chat-server/libs/go-restful"
	"github.com/spring1843/chat-server/src/shared/rest"
)

type MockedReaderCloser struct {
	io.Reader
}

func (MockedReaderCloser) Close() error { return nil }

func TestCanDecorateResponseWithError(t *testing.T) {
	resp := new(rest.Resp)
	restError := rest.RespError{
		Severity:             rest.Error,
		HumanFriendlyMessage: `test error`,
		ShortMessage:         `test-error`,
	}
	resp.AddError(restError)

	httpRequest := &http.Request{
		Host: `localhost`,
		URL: &url.URL{
			Host: `localhost`,
		},
	}

	request := restful.NewRequest(httpRequest)
	resp.DecorateResponse(request)

	if resp.Errors[0].ShortMessage != `test-error` {
		t.Error("Error was not added to response")
	}
}

func TestCanParseRequestBody(t *testing.T) {
	type RequestObject struct {
		Foo string `json:"foo"`
		Baz int    `json:"baz"`
	}

	requestObject := new(RequestObject)
	testJson := `{"foo":"test","baz" : 1}`
	requestBody := MockedReaderCloser{bytes.NewBufferString(testJson)}

	httpRequest := &http.Request{
		Host: `localhost`,
		Body: requestBody,
		URL: &url.URL{
			Host: `localhost`,
		},
		ContentLength: int64(len(testJson)),
	}

	restfulRequest := restful.NewRequest(httpRequest)
	err := rest.ParseRequestBody(restfulRequest, requestObject)
	if err != nil {
		t.Errorf("Error failed parsing JSON:%s", err.Error())
	}

	if requestObject.Foo != `test` || requestObject.Baz != 1 {
		t.Errorf("Json request was not parsed properly %v", requestObject)
	}
}
