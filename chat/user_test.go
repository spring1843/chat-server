package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var (
	user1 = chat.NewUser("u1")
	user2 = chat.NewUser("u2")
	user3 = chat.NewUser("u3")
)

func Test_CanIgnore(t *testing.T) {
	user1.Ignore(user2.NickName)
	if user1.HasIgnored(user2.NickName) != true {
		t.Errorf("User was not ignored when he should have been")
	}

	if user1.HasIgnored(user3.NickName) != false {
		t.Errorf("User was ignored when he should not have been")
	}
}

func Test_CanWriteToUser(t *testing.T) {
	fakeWriter := NewMockedChatConnection()
	user1 := chat.NewConnectedUser(server, fakeWriter)

	go user1.SetOutgoing(`foo`)
	msg := user1.GetOutgoing()

	if msg != "foo" {
		t.Errorf("Message was not written to the user. Msg %s", msg)
	}
}

func Test_CanReadFromUser(t *testing.T) {
	fakeReader := NewMockedChatConnection()
	fakeReader.incoming = []byte("foo\n")

	user1 := chat.NewConnectedUser(server, fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
