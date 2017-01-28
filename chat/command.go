package chat

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

func (c *Command) getChatCommand() Command {
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
		commandName := `/` + command.getChatCommand().Name
		if len(input) >= len(commandName) && input[:len(commandName)] == commandName {
			return command, nil
		}
	}
	return nil, errors.New("Command not found")
}

// Execute shows all available chat commands on this server
func (c *HelpCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User param is not set")
	}

	helpMessage := "Here is the list of all available commands\n"
	for _, command := range AllChatCommands {
		helpMessage += command.getChatCommand().Syntax + "\t" + command.getChatCommand().Description + ".\n"
	}

	params.user1.SetOutgoing(helpMessage)
	return nil
}

// Execute lists all the users in a channel to a user
func (c *ListCommand) Execute(params CommandParams) error {
	listMessage := "Here is the list of all users in this channel\n"

	if params.user1 == nil {
		return errors.New("User param is not set")
	}

	if params.user1.Channel == nil {
		return errors.New("User is not in a channel")
	}

	for nickName := range params.user1.Channel.Users {
		if nickName == params.user1.NickName {
			continue
		}
		listMessage += "@" + nickName + ".\n"
	}

	params.user1.SetOutgoing(listMessage)

	return nil
}

// Execute allows a user to ignore another user so to suppress all incoming messages from that user
func (c *IgnoreCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User1 param is not set")
	}

	if params.user2 == nil {
		return errors.New("User2 param is not set")
	}

	if params.user2.NickName == params.user1.NickName {
		return errors.New("We don't let one ignore themselves")
	}

	params.user1.Ignore(params.user2.NickName)
	params.user1.SetOutgoing(params.user2.NickName + " is now ignored.")
	return nil
}

// Execute allows a user to join a channel
func (c *JoinCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User1 param is not set")
	}
	if params.channel == nil {
		chatCommand := c.getChatCommand()
		channelName, err := chatCommand.ParseChannelFromInput(params.rawInput)
		if err != nil {
			return errors.New("Could not parse channel name")
		}
		params.server.AddChannel(channelName)
		params.channel, _ = params.server.GetChannel(channelName)
	}

	if params.user1.Channel != nil && params.user1.Channel.Name == params.channel.Name {
		return errors.New("You are already in channel #" + params.channel.Name)
	}

	params.channel.AddUser(params.user1.NickName)
	params.user1.Channel = params.channel

	params.user1.SetOutgoing("There are " + strconv.Itoa(len(params.channel.Users)) + " other users this channel.")
	params.channel.Broadcast(params.server, `@`+params.user1.NickName+` just joined channel #`+params.channel.Name)
	return nil
}

// Execute allows a user to send a private message to another user
func (c *PrivateMessageCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User1 param is not set")
	}

	if params.user2 == nil {
		return errors.New("User2 param is not set")
	}

	if params.user1.Channel != params.user2.Channel {
		return errors.New("Users are not in the same channel")
	}

	if params.user2.HasIgnored(params.user1.NickName) {
		return errors.New("User has ignored the sender")
	}

	params.server.LogPrintf("message \t @%s to @%s message=%s", params.user1.NickName, params.user2.NickName, params.message)

	now := time.Now()
	go params.user2.SetOutgoing(now.Format(time.Kitchen) + ` - *Private from @` + params.user1.NickName + `: ` + params.message)
	return nil
}

// Execute disconnects a user from server
func (c *QuitCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User1 param is not set")
	}

	if err := params.server.RemoveUser(params.user1.NickName); err != nil {
		return errors.New("Could not remove user afeter quit command")
	}
	if params.user1.Channel != nil {
		params.user1.Channel.RemoveUser(params.user1.NickName)
	}
	if err := params.user1.Disconnect(params.server); err!=nil {
		return errors.New("Could not disconnect user")
	}
	return nil
}
