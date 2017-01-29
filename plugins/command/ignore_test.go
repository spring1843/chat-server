package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/plugins/command"
)

func Test_IgnoreCommand(t *testing.T) {
	fakeConnection1 := chat.NewFakeConnection()
	fakeConnection2 := chat.NewFakeConnection()

	server := chat.NewServer()

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user2.SetNickName(`u2`)

	server.AddUser(user1)
	server.AddUser(user2)

	input := `/ignore @u2`
	ignoreCommand, _ := command.GetCommand(input)
	user1.ExecuteCommand(server, input, ignoreCommand)
	if user1.HasIgnored(user2.GetNickName()) != true {
		t.Errorf("User was not ignored after ignore command executed")
	}
}
