package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var (
	server = chat.NewServer()
)

func Test_CanAddUser(t *testing.T) {
	server.AddUser(user1)
	if !server.IsUserConnected(`u1`) {
		t.Errorf("User is not connected when should have been connected")
	}
	if server.IsUserConnected(`bar`) {
		t.Errorf("User is connected when should not have been connected")
	}
}

func Test_CanRemoveUser(t *testing.T) {
	server.AddUser(user1)
	server.AddUser(user2)

	server.RemoveUser(user1.GetNickName())

	if server.IsUserConnected(`u1`) {
		t.Errorf("User is was not removed when should have been")
	}

	if server.ConnectedUsersCount() != 1 {
		t.Errorf("After adding two users and removing one user total users does not equal 1")
	}
}

func Test_AddChannel(t *testing.T) {
	server.AddChannel(`foo`)

	if server.GetChannelCount() != 1 {
		t.Errorf("Couldn't add a channel")
	}
}

func Test_GetSameChannel(t *testing.T) {
	server.AddChannel(`foo`)
	sameChannel, err := server.GetChannel(`foo`)

	if err != nil || "foo" != sameChannel.GetName() {
		t.Errorf("Couldn't add and get channel")
	}
}
