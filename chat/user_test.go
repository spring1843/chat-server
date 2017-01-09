package chat_test

import (
	"reflect"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_CanIgnore(t *testing.T) {
	user1 := &chat.User{NickName: `u1`}
	user2 := &chat.User{NickName: `u2`}
	user3 := &chat.User{NickName: `u3`}

	user1.Ignore(*user2)
	if user1.HasIgnored(*user2) != true {
		t.Errorf("User was not ignored when he should have been")
	}

	if user1.HasIgnored(*user3) != false {
		t.Errorf("User was ignored when he should not have been")
	}
}

func Test_CanWriteToUser(t *testing.T) {
	fakeWriter := NewMockedChatConnection()
	user1 := chat.NewUser(fakeWriter)

	user1.Outgoing <- `foo`

	if reflect.DeepEqual(fakeWriter.outgoing, []byte("foo\n")) == false {
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
