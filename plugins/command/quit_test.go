package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
	"github.com/spring1843/chat-server/plugins/command"
)

func Test_QuitCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	user.SetNickName(`foo`)
	server.AddUser(user)

	if server.IsUserConnected(`foo`) != true {
		t.Errorf("User was  disconnected without runnign the quit command")
	}

	input := `/quit`
	quitCommand, _ := command.GetCommand(input)
	user.ExecuteCommand(server, input, quitCommand)

	if server.IsUserConnected(`foo`) {
		t.Errorf("User was not disconnected after running quit command")
	}
}
