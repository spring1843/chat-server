package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestMessageCommand(t *testing.T) {
	server := chat.NewServer()
	user1, user2 := chat.NewUser("u1"), chat.NewUser("u2")

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.SetName(`r`)
	user1.SetChannel(channel.Name)
	user2.SetChannel(channel.Name)

	msg := "foo"
	input := `/msg @u2 ` + msg
	if _, err := user1.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}

	incoming := user2.GetOutgoing()
	if !strings.Contains(incoming, msg) {
		t.Errorf("Message was not read from the user, expected %s got %s", msg, incoming)
	}
}
