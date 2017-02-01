package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestUserCanGettersAndSetters(t *testing.T) {

	user1 := chat.NewUser("u1")
	user2 := chat.NewUser("u2")
	user3 := chat.NewUser("u3")

	user1.Ignore(user2.GetNickName())
	if user1.HasIgnored(user2.GetNickName()) != true {
		t.Errorf("User was not ignored when he should have been")
	}

	if user1.HasIgnored(user3.GetNickName()) != false {
		t.Errorf("User was ignored when he should not have been")
	}

	if user1.SetNickName("nick"); user1.GetNickName() != "nick" {
		t.Errorf("User nickname was not set properly")
	}

	if user1.SetChannel("baz"); user1.GetChannel() != "baz" {
		t.Errorf("User nickname was not set properly")
	}
}
