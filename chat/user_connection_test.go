package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_CanWriteToUser(t *testing.T) {
	fakeWriter := chat.NewMockedChatConnection()
	user1 := chat.NewConnectedUser(server, fakeWriter)

	go user1.SetOutgoing(`foo`)
	msg := user1.GetOutgoing()

	if msg != "foo" {
		t.Errorf("Message was not written to the user. Msg %s", msg)
	}
}

func Test_CanReadFromUser(t *testing.T) {
	fakeReader := chat.NewMockedChatConnection()
	fakeReader.Incoming = []byte("foo\n")

	user1 := chat.NewConnectedUser(server, fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
