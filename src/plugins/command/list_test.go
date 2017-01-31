package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestListCommand(t *testing.T) {
	fakeConnection1 := fake.NewFakeConnection()
	fakeConnection2 := fake.NewFakeConnection()
	fakeConnection3 := fake.NewFakeConnection()

	server := chat.NewServer()
	server.AddChannel(`foo`)

	user1 := chat.NewConnectedUser(fakeConnection1)
	user1.SetNickName(`u1`)
	user1.Listen(server)

	user2 := chat.NewConnectedUser(fakeConnection2)
	user2.SetNickName(`u2`)
	user2.Listen(server)

	user3 := chat.NewConnectedUser(fakeConnection3)
	user3.SetNickName(`u3`)
	user3.Listen(server)

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

	msg, err := fakeConnection1.ReadString(255)
	if err != nil {
		t.Fatalf("Error reading from connection. Error %s", err)
	}
	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Fatalf("List command did not show all users in the room")
	}
}
