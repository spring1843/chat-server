package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func Test_CanWriteToUser(t *testing.T) {
	fakeWriter := fake.NewFakeConnection()
	user1 := chat.NewConnectedUser(server, fakeWriter)

	go user1.SetOutgoing(`foo`)
	msg := user1.GetOutgoing()

	if msg != "foo" {
		t.Errorf("Message was not written to the user. Msg %s", msg)
	}
}

func Test_CanReadFromUser(t *testing.T) {
	fakeReader := fake.NewFakeConnection()
	fakeReader.Incoming = []byte("foo\n")

	user1 := chat.NewConnectedUser(server, fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
