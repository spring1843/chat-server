package chat

import (
	"errors"
	"strconv"
	"time"
)

// Chat commands are executed by a user on a server
type ChatCommand struct {
	Name           string
	Syntax         string
	Description    string
	RequiredParams []string
}

func (c *ChatCommand) getChatCommand() ChatCommand {
	return *c
}

// Every command must be executable
type Executable interface {
	Execute(params CommandParams) error
	getChatCommand() ChatCommand
	ParseNickNameFomInput(input string) (string, error)
	ParseChannelFromInput(input string) (string, error)
	ParseMessageFromInput(input string) (string, error)
	ParseCommandFromInput(input string) (string, error)
}

// Variable containing all valid chat commands supported by this server
var AllChatCommands []Executable

func init() {
	AllChatCommands = []Executable{
		&HelpCommand{
			ChatCommand{
				Name:           `help`,
				Syntax:         `/help`,
				Description:    `Shows the list of all available commands`,
				RequiredParams: []string{`user1`},
			}},
		&ListCommand{
			ChatCommand{
				Name:           `list`,
				Syntax:         `/list`,
				Description:    `Lists user nicknames in the current channel`,
				RequiredParams: []string{`user1`},
			}},
		&IgnoreCommand{
			ChatCommand{
				Name:           `ignore`,
				Syntax:         `/ignore @nickname`,
				Description:    `Ignore a user, followed by user nickname. An ignored user can not send you private messages`,
				RequiredParams: []string{`user1`, `user2`},
			}},
		&JoinCommand{
			ChatCommand{
				Name:           `join`,
				Syntax:         `/join #channel`,
				Description:    `Join a channel`,
				RequiredParams: []string{`user1`, `channel`},
			}},
		&PrivateMessageCommand{
			ChatCommand{
				Name:           `msg`,
				Syntax:         `/msg @nickname message`,
				Description:    `Send a private message to a user in the same channel`,
				RequiredParams: []string{`user1`, `user2`, `message`},
			}},
		&QuitCommand{
			ChatCommand{
				Name:           `quit`,
				Syntax:         `/quit`,
				Description:    `Quit chat server`,
				RequiredParams: []string{`user1`},
			}},
	}
}

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

type HelpCommand struct {
	ChatCommand
}

// Shows all available chat commands on this server
func (c *HelpCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User param is not set")
	}

	helpMessage := "Here is the list of all available commands\n"
	for _, command := range AllChatCommands {
		helpMessage += command.getChatCommand().Syntax + "\t" + command.getChatCommand().Description + ".\n"
	}

	params.user1.Outgoing <- helpMessage
	return nil
}

type ListCommand struct {
	ChatCommand
}

// Lists all the users in a channel to a user
func (c *ListCommand) Execute(params CommandParams) error {
	listMessage := "Here is the list of all users in this channel\n"

	if params.user1 == nil {
		return errors.New("User param is not set")
	}

	if params.user1.Channel == nil {
		return errors.New("User is not in a channel")
	}

	for _, user := range params.user1.Channel.Users {
		if user.NickName == params.user1.NickName {
			continue
		}
		listMessage += "@" + user.NickName + ".\n"
	}

	params.user1.Outgoing <- listMessage

	return nil
}

type IgnoreCommand struct {
	ChatCommand
}

// Allows a user to ignore another user so to suppress all incoming messages from that user
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

	params.user1.Ignore(*params.user2)
	params.user1.Outgoing <- params.user2.NickName + " is now ignored."
	return nil
}

type JoinCommand struct {
	ChatCommand
}

// Allows a user to join a channel
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
		params.user1.Server.AddChannel(channelName)
		params.channel, _ = params.user1.Server.GetChannel(channelName)
	}

	if params.user1.Channel != nil && params.user1.Channel.Name == params.channel.Name {
		return errors.New("You are already in channel #" + params.channel.Name)
	}

	params.channel.AddUser(params.user1)
	params.user1.Channel = params.channel

	params.user1.Outgoing <- "There are " + strconv.Itoa(len(params.channel.Users)) + " other users this channel."

	params.channel.Broadcast(`@` + params.user1.NickName + ` just joined channel #` + params.channel.Name)
	return nil
}

type PrivateMessageCommand struct {
	ChatCommand
}

// Allows a user to send a private message to another user
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

	if params.user2.HasIgnored(*params.user1) {
		return errors.New("User has ignored the sender")
	}

	params.user1.Server.LogPrintf("message \t @%s to @%s message=%s", params.user1.NickName, params.user2.NickName, params.message)

	now := time.Now()
	params.user2.Outgoing <- now.Format(time.Kitchen) + ` - *Private from @` + params.user1.NickName + `: ` + params.message

	return nil
}

type QuitCommand struct {
	ChatCommand
}

// Disconnects a user from server
func (c *QuitCommand) Execute(params CommandParams) error {
	if params.user1 == nil {
		return errors.New("User1 param is not set")
	}

	params.user1.Server.RemoveUser(params.user1)
	if params.user1.Channel != nil {
		params.user1.Channel.RemoveUser(params.user1)
	}
	params.user1.Disconnect()

	return nil
}
