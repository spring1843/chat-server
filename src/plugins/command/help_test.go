package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestHelpCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()

	input := "/help\n"
	n, err := fakeConnection.WriteString(input)
	if err != nil {
		t.Fatalf("Failed writing to connection. Error %s", err)
	}
	if n != len(input) {
		t.Fatalf("Wrong length after write. Expected %d, got %d.", len(input), n)
	}

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	user.SetNickName("foo")
	server.AddUser(user)

	chat.ExpectOutgoing(t, user, 5, "Shows the list of all available commands")
}
