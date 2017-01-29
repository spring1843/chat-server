package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/plugins/command"
)

func Test_MessageCommand(t *testing.T) {
	server := chat.NewServer()

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.Name = `r`
	user1.SetChannel(channel.Name)
	user2.SetChannel(channel.Name)
	input := `/msg @u2 foo`
	messageCommand, err := command.GetCommand(input)
	if err != nil {
		t.Fatalf("Failed getting message command. Error %s", err)
	}

	if err := user1.ExecuteCommand(server, input, messageCommand); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}
	msg := user2.GetOutgoing()

	if strings.Contains(msg, "- *Private from @u1: foo") != true {
		t.Errorf("Private message was not received. Last message %s", msg)
	}
}
