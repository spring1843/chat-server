package rest_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/rest"
)

func TestCanStartAndGetStatus(t *testing.T) {
	config := config.Config{
		IP:       `0.0.0.0`,
		RestPort: 4001,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	restfulAPI := rest.NewRESTfulAPI(config, chatServer)
	go restfulAPI.ListenAndServe()

	response, err := http.Get("http://localhost:4001/status")
	if err != nil {
		t.Errorf("Error making http call to status")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Errorf("Error reading response")
	}

	contentsAsString := string(contents)
	if strings.Contains(contentsAsString, `health`) == false {
		t.Errorf("Status response did not contain health information")
	}
}
