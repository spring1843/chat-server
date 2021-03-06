package chat

import (
	"time"

	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/plugins/command"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

var errIncomingTimeOut = errs.Newf("Read timeout after %s", "15s")

// GetIncoming gets the incoming message from the user
func (u *User) GetIncoming() (string, error) {
	var incoming string
	select {
	case <-time.After(time.Second * 15):
		return incoming, errIncomingTimeOut
	case incoming = <-u.incoming:
	}
	return incoming, nil
}

// SetIncoming sets an incoming message from the user
func (u *User) SetIncoming(message string) {
	u.incoming <- message
}

// HandleNewInput interprets user input and lets chatServer handle it
func (u *User) HandleNewInput(chatServer *Server, userInput string) (bool, error) {
	if u.GetNickName() == "" {
		// This is from a user who is not identified yet, we do not do anything
		// about his input
		logs.Infof("Unidentified user sent input %q", userInput)
		u.SetIncoming(userInput)
		return true, nil
	}

	if command.IsInputExecutable(userInput) {
		return u.handleCommandInput(chatServer, userInput)
	}

	// If it's not a command it's a chat message to broadcast into the channel
	if u.GetChannel() == "" {
		u.SetOutgoing(plugins.UserOutPutTypeLogErr, "You need to join a channel, use /join #channel or use /help for more info.")
		return false, errs.New("User is not in a channel, and input is not a command")
	}
	return u.handleBroadCastInput(chatServer, userInput)
}

func (u *User) handleBroadCastInput(chatServer *Server, userInput string) (bool, error) {
	channel, err := chatServer.GetChannel(u.GetChannel())
	if err != nil {
		return false, errs.Wrap(err, "Error getting channel from server")
	}
	channel.Broadcast(chatServer, `@`+u.GetNickName()+`: `+userInput)
	return true, nil
}

func (u *User) handleCommandInput(chatServer *Server, input string) (bool, error) {
	userCommand, err := command.FromString(input)
	if err != nil {
		u.SetOutgoingf(plugins.UserOutPutTypeLogErr, "Invalid command, use /help for more info. Error: %s", err.Error())
		logs.ErrIfErrf(err, "Failed executing %s command by @%s", input, u.GetNickName())
		return false, errs.Wrap(err, "Error getting command from user input.")
	}

	commandParams, err := u.GetCommandParams(chatServer, input, userCommand)
	if err != nil {
		u.SetOutgoingf(plugins.UserOutPutTypeLogErr, "Error executing your command. %s", err)
		return false, errs.Wrap(err, "Couldn't get command params")
	}

	if err = userCommand.Execute(*commandParams); err != nil {
		logs.ErrIfErrf(err, "error \t @%s command=%s error=%s", u.GetNickName(), input, err)
		return false, errs.Wrapf(err, "Couldn't execute command %q.", input)
	}

	logs.Infof("User @%s executed command: %s", u.GetNickName(), input)
	return true, nil
}

// GetCommandParams looks at command parameters in userInput and populates the parameters for command execution
func (u *User) GetCommandParams(chatServer *Server, userInput string, executable command.Executable) (*command.Params, error) {
	commandParams := &command.Params{
		User1:    u,
		RawInput: userInput,
		Server:   chatServer,
	}

	if executable.RequiresParam(`user2`) {
		nickname, err := command.ParseNickNameFomInput(userInput)
		if err != nil {
			return nil, errs.Wrap(err, "Could not find the required @nickname in the input")
		}

		user2, err := chatServer.GetUser(nickname)
		if err != nil {
			return nil, errs.Wrap(err, "User "+nickname+" + is not connected to this server")
		}
		commandParams.User2 = user2
	}

	if executable.RequiresParam(`channel`) {
		channelName, err := command.ParseChannelFromInput(userInput)
		if err != nil {
			return nil, errs.Wrapf(err, "Could not find the required #channel in the input. Input %s", userInput)
		}
		channel, err := chatServer.GetChannel(channelName)
		if err != nil {
			// If channels don't exist on join command we create one so that channels can start
			switch executable.GetChatCommand().Name {
			case "join":
				chatServer.AddChannel(channelName)
				channel, err = chatServer.GetChannel(channelName)
				if err != nil {
					logs.ErrIfErrf(err, "Couldn't get the channel just created #%s", channelName)
				}
			default:
				return nil, errs.Wrapf(err, "Could not get channel from the server. Channel Name %s", channelName)
			}
		}
		commandParams.Channel = channel
	}

	if executable.RequiresParam(`message`) {
		message, err := command.ParseMessageFromInput(userInput)
		if err != nil {
			return nil, errs.Wrap(err, "Could not required message in the input")
		}
		commandParams.Message = message
	}
	return commandParams, nil
}
