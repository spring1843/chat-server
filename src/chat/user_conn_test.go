package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestCanWriteToUser(t *testing.T) {
	user1 := chat.NewUser("bar")

	msg := "foo"
	go user1.SetOutgoing(msg)

	outgoing := user1.GetOutgoing()
	if outgoing != msg {
		t.Errorf("Received message %q which is not equal to %q", outgoing, msg)

	}
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

	user1 := chat.NewConnectedUser(fakeReader)
	msg := user1.GetIncoming()

	if msg != "foo" {
		t.Errorf("Message was not read from the user, got %s", msg)
	}
}
