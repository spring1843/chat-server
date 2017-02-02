package chat_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestHandleCommandInput(t *testing.T) {
	server := chat.NewServer()
	user1 := chat.NewUser("u1")

	server.AddUser(user1)

	go func(t *testing.T) {
		input := `/join #p`
		if _, err := user1.HandleNewInput(server, input); err != nil {
			t.Fatalf("Failed executing help command. Error %s", err)
		}
	}(t)

	incoming := user1.GetOutgoing()
	if !strings.Contains(incoming, "You are the first") {
		t.Errorf("Message was not sent to the user, expected channel welcome message to be part of %s", incoming)
	}
}

func TestHandleCommandInputFailure(t *testing.T) {
	server := chat.NewServer()
	user1 := chat.NewUser("u1")

	server.AddUser(user1)

	go func(t *testing.T) {
		input := `/yelp`
		if _, err := user1.HandleNewInput(server, input); err == nil {
			t.Fatal("Did not fail executing invalid command", err)
		}
	}(t)

	incoming := user1.GetOutgoing()
	if !strings.Contains(incoming, "not found") {
		t.Errorf("Message was not read from the user, expected quit to be part of %s", incoming)
	}

	go func(t *testing.T) {
		input := `/help`
		if _, err := user1.HandleNewInput(server, input); err != nil {
			t.Fatalf("Failed executing help command. Error %s", err)
		}
	}(t)

	incoming = user1.GetOutgoing()
	if !strings.Contains(incoming, "quit") {
		t.Errorf("Message was not sent to the user, expected quit to be part of %s", incoming)
	}
}

func TestHandleBroadCastInput(t *testing.T) {
	server := chat.NewServer()
	user1, user2 := chat.NewUser("u1"), chat.NewUser("u2")

	channelName := "bar"
	server.AddChannel(channelName)
	channel, err := server.GetChannel(channelName)
	if err != nil {
		t.Fatalf("Error getting channel just added. %s", err)
	}

	server.AddUser(user1)
	server.AddUser(user2)

	channel.AddUser("u1")
	channel.AddUser("u2")

	user1.SetChannel(channelName)
	user2.SetChannel(channelName)

	go func(t *testing.T) {
		input := `foo`
		if _, err := user1.HandleNewInput(server, input); err != nil {
			t.Fatalf("Failed executing help command. Error %s", err)
		}
	}(t)

	incoming := user2.GetOutgoing()
	if !strings.Contains(incoming, "foo") {
		t.Errorf("Message was not read from the user, expected quit to be part of %s", incoming)
	}
}
