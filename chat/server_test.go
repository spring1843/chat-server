package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var (
	server = chat.NewServer()
)

func TestCanAddUser(t *testing.T) {
	server.AddUser(user1)

	user, err := server.GetUser(user1.GetNickName())
	if err != nil {
		t.Errorf("Couldn't get the user just added. Error %s", err)
	}

	if !server.IsUserConnected(user.GetNickName()) {
		t.Errorf("User is not connected when should have been connected")
	}
	if server.IsUserConnected(`bar`) {
		t.Errorf("User is connected when should not have been connected")
	}
}

func TestCanRemoveUser(t *testing.T) {
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

func TestAddChannel(t *testing.T) {
	server.AddChannel(`foo`)

	if server.GetChannelCount() != 1 {
		t.Errorf("Couldn't add a channel")
	}
}

func TestGetSameChannel(t *testing.T) {
	server.AddChannel(`foo`)
	sameChannel, err := server.GetChannel(`foo`)

	if err != nil || "foo" != sameChannel.GetName() {
		t.Errorf("Couldn't add and get channel")
	}
}
