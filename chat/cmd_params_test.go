package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_ParsingChatCommands(t *testing.T) {
	cmd1 := `/msg @nickname #channel foo bar baz`
	cmd2 := `#channel @nickname foo bar baz /foo`
	cmd3 := `foo bar baz @nickname #channel`
	cmd4 := `foo bar baz`
	cmd5 := `#channel`

	chatCommand := new(chat.ChatCommand)

	if output, _ := chatCommand.ParseCommandFromInput(cmd1); output != `msg` {
		t.Errorf("Could not parse command name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseCommandFromInput(cmd2); output != `foo` {
		t.Errorf("Could not parse command name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseCommandFromInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no command, got %s", output)
	}

	if output, _ := chatCommand.ParseChannelFromInput(cmd1); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseChannelFromInput(cmd2); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseChannelFromInput(cmd3); output != `channel` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseNickNameFomInput(cmd1); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseNickNameFomInput(cmd2); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseNickNameFomInput(cmd3); output != `nickname` {
		t.Errorf("Could not parse channel name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseMessageFromInput(cmd1); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseMessageFromInput(cmd2); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseMessageFromInput(cmd3); output != `foo bar baz` {
		t.Errorf("Could not parse message name properly, got %s", output)
	}

	if output, _ := chatCommand.ParseNickNameFomInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no nickname, got %s", output)
	}

	if output, _ := chatCommand.ParseChannelFromInput(cmd4); output != `` {
		t.Errorf("Did not parse empty when there's no channel, got %s", output)
	}
	if output, _ := chatCommand.ParseMessageFromInput(cmd5); output != `` {
		t.Errorf("Did not parse empty when there's no message, got %s", output)
	}

}

func Test_DoesCommandRequireParam(t *testing.T) {
	fakeCommand := &chat.QuitCommand{
		ChatCommand: chat.ChatCommand{
			Name:           `quit`,
			Syntax:         `/quit`,
			Description:    `Quit chat server`,
			RequiredParams: []string{`user1`},
		},
	}

	if fakeCommand.ChatCommand.DoesCommandRequireParam(`user1`) == false {
		t.Errorf("Required param user1 was not seen as required.")
	}

	if fakeCommand.ChatCommand.DoesCommandRequireParam(`user2`) == true {
		t.Errorf("Not-required param user2 was not seen as required.")
	}
}
