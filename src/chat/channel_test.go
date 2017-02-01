package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

var channel = chat.NewChannel()

func TestCanAddUsers(t *testing.T) {
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
	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())

	var chatServer = chat.NewServer()
	chatServer.AddUser(user1)
	chatServer.AddUser(user2)

	msg := "foo"
	go channel.Broadcast(chatServer, msg)

	chat.ExpectOutgoing(t, user1, 5, msg)
	chat.ExpectOutgoing(t, user2, 5, msg)
}

func TestChatCanGettersAndSetters(t *testing.T) {
	if channel.SetName("baz"); channel.GetName() != "baz" {
		t.Errorf("Channel name was not set properly")
	}
}
