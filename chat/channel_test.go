package chat_test

import (
	"testing"

	"strings"

	"github.com/spring1843/chat-server/chat"
)

func Test_CanAddUsers(t *testing.T) {
	var (
		channel = chat.NewChannel()
		user1   = new(chat.User)
		user2   = new(chat.User)
	)

	channel.AddUser(user1)
	channel.AddUser(user2)
	if len(channel.Users) != 2 {
		t.Errorf("Users couldn't be added to the channel")
	}
}

func Test_CanBroadCast(t *testing.T) {
	var (
		channel = chat.NewChannel()
		user1   = new(chat.User)
		user2   = new(chat.User)
	)

	user1.Outgoing = make(chan string)
	user2.Outgoing = make(chan string)

	channel.AddUser(user1)
	channel.AddUser(user2)

	go channel.Broadcast(`foo`)

	msg1 := <-user1.Outgoing
	msg2 := <-user2.Outgoing

	if strings.Contains(msg1, `foo`) != true || msg1 != msg2 {
		t.Errorf("Message wasn't broadcasted to the users in the channel")
	}
}
