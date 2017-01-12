package integration_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/integration"
	"github.com/spring1843/chat-server/rest"
)

func TestCanStartAndGetStatus(t *testing.T) {
	config := integration.Config{
		IP:       `0.0.0.0`,
		RestPort: 4001,
		LogFile:  `/dev/null`,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	testFile, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	chatServer.SetLogFile(testFile)

	rest.Start(chatServer, config)

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
