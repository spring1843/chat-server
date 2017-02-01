package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestCanReadWriteToFromUser(t *testing.T) {
	user1 := chat.NewUser("bar")

	input := "foo"
	go user1.SetOutgoing(input)

	outgoing := user1.GetOutgoing()
	if outgoing != input {
		t.Errorf("Received message %q which is not equal to %q", outgoing, input)

	}

	user1 = chat.NewUser("bar")
	go user1.SetIncoming(input)

	incoming := user1.GetIncoming()
	if incoming != input {
		t.Errorf("Message was not read from the user, got %s", incoming)
	}
}
