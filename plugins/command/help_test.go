package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func TestHelpCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()
	fakeConnection.WriteString("/help\n")

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	user.SetNickName("foo")
	server.AddUser(user)
	msg := user.GetOutgoing()

	if strings.Contains(msg, "Shows the list of all available commands") != true {
		t.Errorf("Help command did not output description of help command. Got %s", msg)
	}
}
