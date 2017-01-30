package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/plugins/command"
)

var (
	user1 = chat.NewUser("u1")
	user2 = chat.NewUser("u2")
)

func TestCanGetCommand(t *testing.T) {
	command, err := command.FromString("/join")
	if err != nil {
		t.Errorf("Couldn't get join command. Error: %s", err)
	}
	sameCommand := command.GetChatCommand()
	if sameCommand.Name != "join" {
		t.Fatalf("Couldn't get the same command. Expected /join got %s", sameCommand.Name)
	}

}

func TestCanValidate(t *testing.T) {
	var (
		invalidCommand1 = ``
		invalidCommand2 = `badcommand`
		invalidCommand3 = `/badcommand`
		validCommand1   = `/help`
		validCommand2   = `/join`
	)

	if _, err := command.FromString(invalidCommand1); err == nil {
		t.Errorf("Invalid command was detected valid. command: %s Error: %s", invalidCommand1, err)
	}

	if _, err := command.FromString(invalidCommand2); err == nil {
		t.Errorf("Invalid command was detected valid. command: %s Error: %s", invalidCommand2, err)
	}

	if _, err := command.FromString(invalidCommand3); err == nil {
		t.Errorf("Invalid command was detected valid. command: %s Error: %s", invalidCommand3, err)
	}

	if _, err := command.FromString(validCommand1); err != nil {
		t.Errorf("Valid command was detected invalid. command: %s Error: %s", validCommand1, err)
	}

	if _, err := command.FromString(validCommand2); err != nil {
		t.Errorf("Valid command was detected invalid, command: %s. Error: %s", validCommand2, err)
	}
}
