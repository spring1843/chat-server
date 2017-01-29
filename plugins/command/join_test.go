package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
	"github.com/spring1843/chat-server/plugins/command"
)

func Test_JoinCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()

	server := chat.NewServer()
	user1 := chat.NewConnectedUser(server, fakeConnection)
	user1.SetNickName("u1")
	server.AddUser(user1)

	input := `/join #r`
	joinCommand, err := command.FromString(input)
	if err != nil {
		t.Errorf("Could not get an instance of join command")
	}

	user1.ExecuteCommand(server, input, joinCommand)
	if user1.GetChannel() == "" {
		t.Errorf("User did not join the #r channel when he should have")
	}
}
