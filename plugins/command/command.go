package command

import "github.com/spring1843/chat-server/plugins/errs"

type (
	// Command is executed by a user on a server
	Command struct {
		Name           string
		Syntax         string
		Description    string
		RequiredParams []string
	}
	// HelpCommand Shows list of available commands
	HelpCommand struct {
		Command
	}
	// ListCommand list available users in a channel
	ListCommand struct {
		Command
	}
	// IgnoreCommand allows a user to ignore another user
	IgnoreCommand struct {
		Command
	}
	// JoinCommand allows user to join a channel
	JoinCommand struct {
		Command
	}
	// PrivateMessageCommand allows a channel to privately message a user
	PrivateMessageCommand struct {
		Command
	}
	// QuitCommand allows a user to disconnect from the server
	QuitCommand struct {
		Command
	}
)

// GetChatCommand returns this command
func (c *Command) GetChatCommand() Command {
	return *c
}

// AllChatCommands all valid chat commands supported by this server
var AllChatCommands []Executable

func init() {
	AllChatCommands = []Executable{
		&HelpCommand{
			Command{
				Name:           `help`,
				Syntax:         `/help`,
				Description:    `Shows the list of all available commands`,
				RequiredParams: []string{`user1`},
			}},
		&ListCommand{
			Command{
				Name:           `list`,
				Syntax:         `/list`,
				Description:    `Lists user nicknames in the current channel`,
				RequiredParams: []string{`user1`},
			}},
		&IgnoreCommand{
			Command{
				Name:           `ignore`,
				Syntax:         `/ignore @nickname`,
				Description:    `Ignore a user, followed by user nickname. An ignored user can not send you private messages`,
				RequiredParams: []string{`user1`, `user2`},
			}},
		&JoinCommand{
			Command{
				Name:           `join`,
				Syntax:         `/join #channel`,
				Description:    `Join a channel`,
				RequiredParams: []string{`user1`, `channel`},
			}},
		&PrivateMessageCommand{
			Command{
				Name:           `msg`,
				Syntax:         `/msg @nickname message`,
				Description:    `Send a private message to a user in the same channel`,
				RequiredParams: []string{`user1`, `user2`, `message`},
			}},
		&QuitCommand{
			Command{
				Name:           `quit`,
				Syntax:         `/quit`,
				Description:    `Quit chat server`,
				RequiredParams: []string{`user1`},
			}},
	}
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
