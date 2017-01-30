package command_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/plugins/command"
)

func TestParsingChatCommands(t *testing.T) {
	cmd1 := `/msg @nickname #channel foo bar baz`
	cmd2 := `#channel @nickname foo bar baz /foo`
	cmd3 := `foo bar baz @nickname #channel`
	cmd4 := `foo bar baz`
	cmd5 := `#channel`

	if output, _ := command.ParseCommandFromInput(cmd1); output != `msg` {
		t.Errorf("Could not parse command name properly, got %s", output)
	}

	if output, _ := command.ParseCommandFromInput(cmd2); output != `foo` {
		t.Errorf("Could not parse command name properly, got %s", output)
	}

	if output, _ := command.ParseCommandFromInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no command, got %s", output)
	}

	if output, _ := command.ParseChannelFromInput(cmd1); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseChannelFromInput(cmd2); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseChannelFromInput(cmd3); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseNickNameFomInput(cmd1); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseNickNameFomInput(cmd2); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseNickNameFomInput(cmd3); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := command.ParseMessageFromInput(cmd1); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := command.ParseMessageFromInput(cmd2); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := command.ParseMessageFromInput(cmd3); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := command.ParseNickNameFomInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no nickname, got %s", output)
	}

	if output, _ := command.ParseChannelFromInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no channel, got %s", output)
	}
	if output, _ := command.ParseMessageFromInput(cmd5); output != `` {
		t.Errorf("Did not parse empty when there's no message, got %s", output)
	}

}

func TestDoesCommandRequireParam(t *testing.T) {
	fakeCommand := &command.QuitCommand{
		Command: command.Command{
			Name:           `quit`,
			Syntax:         `/quit`,
			Description:    `Quit chat server`,
			RequiredParams: []string{`user1`},
		},
	}

	if fakeCommand.Command.RequiresParam(`user1`) == false {
		t.Errorf("Required param user1 was not seen as required.")
	}

	if fakeCommand.Command.RequiresParam(`user2`) == true {
		t.Errorf("Not-required param user2 was not seen as required.")
	}
}
