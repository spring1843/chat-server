package chat_test

import (
	"testing"

	"strings"

	"github.com/spring1843/chat-server/chat"
)

var channel = &chat.Channel{Name: "foo", Users: make(map[string]bool)}

func Test_CanAddUsers(t *testing.T) {
	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())
	if len(channel.Users) != 2 {
		t.Errorf("Users couldn't be added to the channel")
	}
}

func Test_CanBroadCast(t *testing.T) {
	channel.AddUser(user1.GetNickName())
	channel.AddUser(user2.GetNickName())

	var chatServer = chat.NewServer()
	chatServer.AddUser(user1)
	chatServer.AddUser(user2)

	go channel.Broadcast(chatServer, `foo`)

	msg1 := user1.GetOutgoing()
	msg2 := user2.GetOutgoing()

	if strings.Contains(msg1, `foo`) != true || msg1 != msg2 {
		t.Errorf("Message wasn't broadcasted to the users in the channel")
	}
}
