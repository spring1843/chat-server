package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func TestCanWriteToUser(t *testing.T) {
	fakeWriter := fake.NewFakeConnection()
	user1 := chat.NewConnectedUser(server, fakeWriter)

	go user1.SetOutgoing(`foo`)
	chat.ExpectOutgoing(t, user1, 5, "foo")
}

func TestCanReadFromUser(t *testing.T) {
	t.Skipf("Racy")
	fakeReader := fake.NewFakeConnection()
	input := "foo\n"
	n, err := fakeReader.WriteString(input)
	if err != nil {
		t.Fatalf("Failed writing to connection. Error %s", err)
	}
	if n != len(input) {
		t.Fatalf("Wrong length after write. Expected %d, got %d.", len(input), n)
	}

	user1 := chat.NewConnectedUser(server, fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
