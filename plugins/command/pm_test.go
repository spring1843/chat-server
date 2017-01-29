package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
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
	if _, err := user1.HandleNewInput(server, input); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}
	msg := user2.GetOutgoing()

	if strings.Contains(msg, "- *Private from @u1: foo") != true {
		t.Errorf("Private message was not received. Last message %s", msg)
	}
}
