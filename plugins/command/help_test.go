package command

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_HelpCommand(t *testing.T) {
	fakeConnection := chat.NewFakeConnection()
	fakeConnection.Incoming = []byte("/help\n")

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	server.AddUser(user)
	msg := user.GetOutgoing()

	if strings.Contains(msg, "Shows the list of all available commands") != true {
		t.Errorf("Help command did not output description of help command")
	}
}
