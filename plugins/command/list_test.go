package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
	"github.com/spring1843/chat-server/plugins/command"
)

func Test_ListCommand(t *testing.T) {
	fakeConnection1 := fake.NewFakeConnection()
	fakeConnection2 := fake.NewFakeConnection()
	fakeConnection3 := fake.NewFakeConnection()

	server := chat.NewServer()
	server.AddChannel(`foo`)

	input := `/join #foo`
	joinCommand, err := command.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user1.SetNickName(`u1`)

	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user1.SetNickName(`u2`)

	user3 := chat.NewConnectedUser(server, fakeConnection3)
	user1.SetNickName(`u3`)

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	user1.ExecuteCommand(server, input, joinCommand)
	user2.ExecuteCommand(server, input, joinCommand)
	user3.ExecuteCommand(server, input, joinCommand)

	input = "/list \n"
	listCommand, err := command.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}
	user1.ExecuteCommand(server, input, listCommand)
	msg := user1.GetOutgoing()

	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}
