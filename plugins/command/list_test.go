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
	user2.SetNickName(`u2`)

	user3 := chat.NewConnectedUser(server, fakeConnection3)
	user3.SetNickName(`u3`)

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	input := `/join #foo`
	if _, err := user1.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed joining. Error: %s", err)
	}
	if _, err := user2.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed joining. Error: %s", err)
	}
	if _, err := user3.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed joining. Error: %s", err)
	}

	if _, err := user1.HandleNewInput(server, "/list"); err != nil {
		t.Fatalf("Failed running list command. Error: %s", err)
	}

	msg := string(fakeConnection1.ReadOutgoing())
	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}
