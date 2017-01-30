package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestMessageCommand(t *testing.T) {
	server := chat.NewServer()

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.SetName(`r`)
	user1.SetChannel(channel.Name)
	user2.SetChannel(channel.Name)

	input := `/msg @u2 foo`
	if _, err := user1.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}
	chat.ExpectOutgoing(t, user2, 5, "- *Private from @u1: foo")
}
