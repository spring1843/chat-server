package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var channel = chat.NewChannel()

func TestCanAddUsers(t *testing.T) {
	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())
	if len(channel.Users) != 2 {
		t.Errorf("Users couldn't be added to the channel")
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
