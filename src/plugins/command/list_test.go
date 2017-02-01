package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestListCommand(t *testing.T) {
	server := chat.NewServer()
	user1, user2 := chat.NewUser("u1"), chat.NewUser("u2")

	server.AddUser(user1)
	server.AddUser(user2)

	channelName := "r"
	server.AddChannel(channelName)
	channel, err := server.GetChannel(channelName)
	if err != nil {
		t.Fatal("Couldn't get the channel just added")
	}

	channel.AddUser("u1")
	channel.AddUser("u2")

	user1.SetChannel(channelName)
	user2.SetChannel(channelName)

	go func(t *testing.T) {
		if _, err := user1.HandleNewInput(server, `/list`); err != nil {
			t.Fatalf("Failed executing message. Error %s", err)
		}
	}(t)

	incoming := user1.GetOutgoing()
	if !strings.Contains(incoming, `u2`) {
		t.Errorf("Message was not read from the user, expected %s got %s", `u2`, incoming)
	}
}
