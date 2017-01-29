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
func FromString(commandName string) (Executable, error) {
	if !IsInputExecutable(commandName) {
		return nil, errs.Wrapf(ErrNotACommand, "%s was not found in available comands.", commandName)
	}

	if command, ok := AllChatCommands[commandName]; ok {
		return command, nil
	}
	return nil, errs.Wrapf(ErrCommadNotFound, "%s was not found in available comands.", commandName)
}
