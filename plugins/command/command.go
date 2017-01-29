package command

import "github.com/spring1843/chat-server/plugins/errs"

// Command is executed by a user on a server
type Command struct {
	Name           string
	Syntax         string
	Description    string
	RequiredParams []string
}

// AllChatCommands all valid chat commands supported by this server
var AllChatCommands = []Executable{
	helpCommand,
	listCommand,
	ignoreCommand,
	joinCommand,
	privateMessageCommand,
	quitCommand,
}

// GetChatCommand returns this command
func (c *Command) GetChatCommand() Command {
	return *c
}

// GetCommand gets a command if it exists
func GetCommand(input string) (Executable, error) {
	if len(input) < 2 || input[0:1] != `/` {
		return nil, errs.New("Input too short to be a command")
	}

	validCommands := AllChatCommands
	for _, command := range validCommands {
		commandName := `/` + command.GetChatCommand().Name
		if len(input) >= len(commandName) && input[:len(commandName)] == commandName {
			return command, nil
		}
	}
	return nil, errs.New("Command not found")
}
