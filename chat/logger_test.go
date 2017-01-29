package chat_test

import (
	"strings"
	"testing"

	"github.com/spring1843/legacy-03/chat-server/chat"
)

func Test_CanLogToFile(t *testing.T) {
	fakeWriter := chat.NewMockedChatConnection()
	server.SetLogFile(fakeWriter)
	server.LogPrintf("test \t foo\n")

	logMessage := string(fakeWriter.Outgoing)

	if strings.Contains(logMessage, `foo`) == false {
		t.Errorf("Did not send log to file")
	}
}
