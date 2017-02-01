package chat_test

import (
	"testing"

	"strings"

	"github.com/spring1843/chat-server/src/chat"
)

var channel = chat.NewChannel()

func TestCanAddUsers(t *testing.T) {
	user1 := chat.NewUser("user1")
	user2 := chat.NewUser("user2")
	channel := chat.NewChannel()

	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())
	if channel.GetUserCount() != 2 {
		t.Errorf("Users couldn't be added to the channel")
	}

	if len(channel.GetUsers()) != 2 {
		t.Errorf("Couldn't get the nicknames of users just added")
	}

	channel.RemoveUser(user1.GetNickName())
	if channel.GetUserCount() != 1 {
		t.Errorf("After removing, user count did not reduce")
	}
}

func TestCanBroadCast(t *testing.T) {
	server := chat.NewServer()
	channel := chat.NewChannel()

	user1 := chat.NewUser("user1")
	user2 := chat.NewUser("user2")

	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())

	server.AddUser(user1)
	server.AddUser(user2)

	msg := "foo"
	channel.Broadcast(server, msg)

	incoming := user1.GetOutgoing()
	if !strings.Contains(incoming, msg) {
		t.Errorf("Message was not read from the user, expected %s got %s", msg, incoming)
	}

	incoming = user2.GetOutgoing()
	if !strings.Contains(incoming, msg) {
		t.Errorf("Message was not read from the user, expected %s got %s", msg, incoming)
	}

}

func TestChatCanGettersAndSetters(t *testing.T) {
	channel := chat.NewChannel()
	if channel.SetName("baz"); channel.GetName() != "baz" {
		t.Errorf("Channel name was not set properly")
	}
}
