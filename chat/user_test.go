package chat_test

import (
	"reflect"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var (
	user1 = &chat.User{NickName: `u1`, Outgoing: make(chan string), IgnoreList: make(map[string]bool)}
	user2 = &chat.User{NickName: `u2`, Outgoing: make(chan string), IgnoreList: make(map[string]bool)}
	user3 = &chat.User{NickName: `u3`, Outgoing: make(chan string), IgnoreList: make(map[string]bool)}
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
	user1 := chat.NewUser(fakeWriter)

	user1.Outgoing <- `foo`

	if reflect.DeepEqual(fakeWriter.ReadOutgoing(), []byte("foo\n")) == false {
		t.Errorf("Message was not written to the user")
	}
}

func Test_CanReadFromUser(t *testing.T) {
	fakeReader := NewMockedChatConnection()
	fakeReader.incoming = []byte("foo\n")

	user1 := chat.NewUser(fakeReader)
	msg := <-user1.Incoming

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
