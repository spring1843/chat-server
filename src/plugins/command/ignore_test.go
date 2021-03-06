package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestIgnoreCommand(t *testing.T) {
	fakeConnection1 := fake.NewFakeConnection()
	fakeConnection2 := fake.NewFakeConnection()

	server := chat.NewServer()

	user1 := chat.NewConnectedUser(fakeConnection1)
	user2 := chat.NewConnectedUser(fakeConnection2)

	user1.Listen(server)
	user2.Listen(server)

	user1.SetNickName(`u1`)
	user2.SetNickName(`u2`)

	server.AddUser(user1)
	server.AddUser(user2)

	if _, err := user1.HandleNewInput(server, `/ignore @u2`); err != nil {
		t.Fatalf("Error ignoring. Error %s", err)
	}
	if user1.HasIgnored(user2.GetNickName()) != true {
		t.Errorf("User was not ignored after ignore command executed")
	}
}
