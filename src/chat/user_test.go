package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

var (
	user1 = chat.NewUser("u1")
	user2 = chat.NewUser("u2")
	user3 = chat.NewUser("u3")
)

func TestCanIgnore(t *testing.T) {
	user1.Ignore(user2.GetNickName())
	if user1.HasIgnored(user2.GetNickName()) != true {
		t.Errorf("User was not ignored when he should have been")
	}

	if user1.HasIgnored(user3.GetNickName()) != false {
		t.Errorf("User was ignored when he should not have been")
	}
}
