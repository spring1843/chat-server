package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestQuitCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()

	server := chat.NewServer()
	user := chat.NewConnectedUser(fakeConnection)
	user.Listen(server)

	user.SetNickName(`foo`)

	server.AddUser(user)

	if !server.IsUserConnected(`foo`) {
		t.Errorf("User was  disconnected without runnign the quit command")
	}

	input := `/quit`
	if _, err := user.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed executing command. Error %s", err)
	}

	if server.IsUserConnected(`foo`) {
		t.Errorf("User was not disconnected after running quit command")
	}
}
