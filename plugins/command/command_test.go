package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/plugins/command"
)

var (
	user1 = chat.NewUser("u1")
	user2 = chat.NewUser("u2")
	user3 = chat.NewUser("u3")
)

func Test_CanValidate(t *testing.T) {
	var (
		invalidCommand1 = ``
		invalidCommand2 = `badcommand`
		invalidCommand3 = `/badcommand`
		validCommand1   = `/help`
		validCommand2   = `/join`
	)

	if _, err := command.GetCommand(invalidCommand1); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand1)
	}

	if _, err := command.GetCommand(invalidCommand2); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand2)
	}

	if _, err := command.GetCommand(invalidCommand3); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand3)
	}

	if _, err := command.GetCommand(validCommand1); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand1)
	}

	if _, err := command.GetCommand(validCommand2); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand2)
	}
}
