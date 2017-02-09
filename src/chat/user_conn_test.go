package chat_test

import (
	"testing"

	"strings"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/plugins"
)

func TestCanReadWriteToFromUser(t *testing.T) {
	user1 := chat.NewUser("bar")

	input := "foo"
	go user1.SetOutgoing(plugins.UserOutPutTUserTest, input)

	outgoing := user1.GetOutgoing()
	if !strings.Contains(outgoing, input) {
		t.Errorf("Received message %q which is not equal to %q", outgoing, input)

	}

	user1 = chat.NewUser("bar")
	go user1.SetIncoming(input)

	incoming, err := user1.GetIncoming()
	if err != nil {
		t.Fatalf("Failed getting incoming from user. Error %s", err)
	}
	if incoming != input {
		t.Errorf("Message was not read from the user, got %s", incoming)
	}
}
