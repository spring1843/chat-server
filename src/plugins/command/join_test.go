package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestJoinCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()

	server := chat.NewServer()
	user1 := chat.NewConnectedUser(fakeConnection)
	user1.Listen(server)

	user1.SetNickName("u1")
	server.AddUser(user1)

	server.AddChannel("r")

	if _, err := user1.HandleNewInput(server, "/join #r"); err != nil {
		t.Fatalf("Could not execute join. Error %s", err)
	}
	if user1.GetChannel() == "" {
		t.Errorf("User did not join the #r channel when he should have")
	}
}
