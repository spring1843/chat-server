package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func Test_JoinCommand(t *testing.T) {
	fakeConnection := fake.NewFakeConnection()
	fakeConnection.Incoming = []byte("/join #r\n")

	server := chat.NewServer()
	user1 := chat.NewConnectedUser(server, fakeConnection)
	server.AddUser(user1)

	msg := user1.GetOutgoing()
	if strings.Contains(msg, "other users this channel") != true {
		t.Errorf("User did not receive welcome message after joining channel received %s instead", msg)
	}

	if user1.GetChannel() != "" && user1.GetChannel() != `r` {
		t.Errorf("User did not join the #r channel when he should have")
	}
}
