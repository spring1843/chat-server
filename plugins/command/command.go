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
		`help`:   helpCommand,
		`list`:   listCommand,
		`ignore`: ignoreCommand,
		`join`:   joinCommand,
		`msg`:    privateMessageCommand,
		`quit`:   quitCommand,
	}
	ErrComadNotFound = errs.New("Command not found")
)

// GetChatCommand returns this command
func (c *Command) GetChatCommand() Command {
	return *c
}

// GetCommand gets a command if it exists
func GetCommand(input string) (Executable, error) {
	if command, ok := AllChatCommands[`/`+input]; ok {
		return command, nil
	}
	return nil, ErrComadNotFound
}
