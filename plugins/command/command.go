package command

import (
	"errors"
	"strconv"
	"time"
)

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
		return nil, errors.New("Input too short to be a command")
	}

	validCommands := AllChatCommands
	for _, command := range validCommands {
		commandName := `/` + command.GetChatCommand().Name
		if len(input) >= len(commandName) && input[:len(commandName)] == commandName {
			return command, nil
		}
	}
	return nil, errors.New("Command not found")
}

// Execute shows all available chat commands on this server
func (c *HelpCommand) Execute(params CommandParams) error {
	if params.User1 == nil {
		return errors.New("User param is not set")
	}

	helpMessage := "Here is the list of all available commands\n"
	for _, command := range AllChatCommands {
		helpMessage += command.GetChatCommand().Syntax + "\t" + command.GetChatCommand().Description + ".\n"
	}

	params.User1.SetOutgoing(helpMessage)
	return nil
}

// Execute lists all the users in a channel to a user
func (c *ListCommand) Execute(params CommandParams) error {
	listMessage := "Here is the list of all users in this channel\n"

	if params.User1 == nil {
		return errors.New("User param is not set")
	}

	if params.User1.GetChannel() == "" {
		return errors.New("User is not in a channel")
	}

	for nickName := range params.Channel.GetUsers() {
		if nickName == params.User1.GetNickName() {
			continue
		}
		listMessage += "@" + nickName + ".\n"
	}

	params.User1.SetOutgoing(listMessage)

	return nil
}

// Execute allows a user to ignore another user so to suppress all incoming messages from that user
func (c *IgnoreCommand) Execute(params CommandParams) error {
	if params.User1 == nil {
		return errors.New("User1 param is not set")
	}

	if params.User2 == nil {
		return errors.New("User2 param is not set")
	}

	if params.User2.GetNickName() == params.User1.GetNickName() {
		return errors.New("We don't let one ignore themselves")
	}

	params.User1.Ignore(params.User2.GetNickName())
	params.User1.SetOutgoing(params.User2.GetNickName() + " is now ignored.")
	return nil
}

// Execute allows a user to join a channel
func (c *JoinCommand) Execute(params CommandParams) error {
	if params.User1 == nil {
		return errors.New("User1 param is not set")
	}

	channelName := ""
	if params.Channel == nil {
		chatCommand := c.GetChatCommand()
		channelName, err := chatCommand.ParseChannelFromInput(params.RawInput)
		if err != nil {
			return errors.New("Could not parse channel name")
		}
		params.Server.AddChannel(channelName)
	}

	if params.User1.GetChannel() != "" && params.User1.GetChannel() == params.Channel.GetName() {
		return errors.New("You are already in channel #" + params.Channel.GetName())
	}

	params.Channel.AddUser(params.User1.GetNickName())
	params.User1.SetChannel(channelName)

	params.User1.SetOutgoing("There are " + strconv.Itoa(params.Channel.GetUserCount()) + " other users this channel.")
	return params.Server.BroadcastInChannel(channelName, `@`+params.User1.GetNickName()+` just joined channel #`+params.Channel.GetName())
}

// Execute allows a user to send a private message to another user
func (c *PrivateMessageCommand) Execute(params CommandParams) error {
	if params.User1 == nil {
		return errors.New("User1 param is not set")
	}

	if params.User2 == nil {
		return errors.New("User2 param is not set")
	}

	if params.User1.GetChannel() != params.User2.GetChannel() {
		return errors.New("Users are not in the same channel")
	}

	if params.User2.HasIgnored(params.User1.GetNickName()) {
		return errors.New("User has ignored the sender")
	}

	params.Server.LogPrintf("message \t @%s to @%s message=%s", params.User1.GetNickName(), params.User2.GetNickName(), params.Message)

	now := time.Now()
	go params.User2.SetOutgoing(now.Format(time.Kitchen) + ` - *Private from @` + params.User1.GetNickName() + `: ` + params.Message)
	return nil
}

// Execute disconnects a user from server
func (c *QuitCommand) Execute(params CommandParams) error {
	if params.User1 == nil {
		return errors.New("User1 param is not set")
	}

	if err := params.Server.RemoveUser(params.User1.GetNickName()); err != nil {
		return errors.New("Could not remove user afeter quit command")
	}
	if params.User1.GetChannel() != "" {
		params.Server.RemoveUserFromChannel(params.User1.GetNickName(), params.User1.GetChannel())
	}

	if err := params.Server.DisconnectUser(params.User1.GetNickName()); err != nil {
		return errors.New("Could not disconnect user")
	}
	return nil
}
