package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func Test_ListCommand(t *testing.T) {
	fakeConnection1 := fake.NewFakeConnection()
	fakeConnection2 := fake.NewFakeConnection()
	fakeConnection3 := fake.NewFakeConnection()

	server := chat.NewServer()
	server.AddChannel(`foo`)

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user1.SetNickName(`u1`)

	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user1.SetNickName(`u2`)

	user3 := chat.NewConnectedUser(server, fakeConnection3)
	user1.SetNickName(`u3`)

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	input := `/join #foo`
	user1.HandleNewInput(server, input)
	user2.HandleNewInput(server, input)
	user3.HandleNewInput(server, input)

	input = "/list \n"
	user1.HandleNewInput(server, input)
	msg := user1.GetOutgoing()

	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}
