package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func Test_IgnoreCommand(t *testing.T) {
	fakeConnection1 := fake.NewFakeConnection()
	fakeConnection2 := fake.NewFakeConnection()

	server := chat.NewServer()

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user2 := chat.NewConnectedUser(server, fakeConnection2)
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
