package command

import (
	"strings"

	"github.com/spring1843/chat-server/plugins/errs"
)

// Command is executed by a user on a server
type Command struct {
	Name           string
	Syntax         string
	Description    string
	RequiredParams []string
}

// AllChatCommands all valid chat commands supported by this server
var (
	AllChatCommands = map[string]Executable{
		`/help`:   helpCommand,
		`/list`:   listCommand,
		`/ignore`: ignoreCommand,
		`/join`:   joinCommand,
		`/msg`:    privateMessageCommand,
		`/quit`:   quitCommand,
	}
	ErrNotACommand    = errs.New("Not a command, commands must start with / and be at least 3 characters")
	ErrCommadNotFound = errs.New("Command not found")
)

// GetChatCommand returns this command
func (c *Command) GetChatCommand() Command {
	return *c
}

// IsInputExecutable checks if a user input is intended to be a command or not
func IsInputExecutable(input string) bool {
	if len(input) > 2 && input[0:1] == "/" {
		return true
	}
	return false
}

// ParseString gets a command if it can find it in a user input string
func FromString(input string) (Executable, error) {
	if !IsInputExecutable(input) {
		return nil, errs.Wrapf(ErrNotACommand, "%s was not found in available comands.", input)
	}

	commandPart := commandPart(input)
	if command, ok := AllChatCommands[commandPart]; ok {
		return command, nil
	}
	return nil, errs.Wrapf(ErrCommadNotFound, "%s was not found in available comands.", commandPart)
}

func commandPart(input string) string {
	return input[:strings.Index(input, " ")]
}
