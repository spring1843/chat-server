package chat_test

import (
	"reflect"
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
	user1 := chat.NewConnectedUser(fakeWriter)

	user1.SetOutgoing(`foo`)

	if reflect.DeepEqual(fakeWriter.ReadOutgoing(), []byte("foo\n")) == false {
		t.Errorf("Message was not written to the user")
	}
}

func Test_CanReadFromUser(t *testing.T) {
	fakeReader := NewMockedChatConnection()
	fakeReader.incoming = []byte("foo\n")

	user1 := chat.NewConnectedUser(fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
