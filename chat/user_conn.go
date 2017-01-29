package chat

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/spring1843/chat-server/drivers"
	"github.com/spring1843/chat-server/plugins/command"
	"github.com/spring1843/chat-server/plugins/errs"
	"github.com/spring1843/chat-server/plugins/logs"
)

const ReadConnectionLimitBytes = 256

// NewConnectedUser returns a new User with a connection
func NewConnectedUser(chatServer *Server, connection drivers.Connection) *User {
	user := NewUser("")
	user.conn = connection
	user.Listen(chatServer)
	return user
}

// Listen starts reading from and writing to a user
func (u *User) Listen(chatServer *Server) {
	go u.ReadFrom(chatServer)
	go u.WriteTo()
}

// GetOutgoing gets the outgoing message for a user
func (u *User) GetOutgoing() string {
	return <-u.outgoing
}

// SetOutgoing sets an outgoing message to the user
func (u *User) SetOutgoing(message string) {
	u.outgoing <- message
}

// GetIncoming gets the incoming message from the user
func (u *User) GetIncoming() string {
	return <-u.incoming
}

// SetIncoming sets an incoming message from the user
func (u *User) SetIncoming(message string) {
	u.incoming <- message
}

// ReadFrom reads data from users and lets chat server interpret it
func (u *User) ReadFrom(chatServer *Server) {
	for {
		message := make([]byte, ReadConnectionLimitBytes)
		if _, err := u.conn.Read(message); err != nil {
			if err == io.EOF {
				continue
			}
			logs.Errf(err, "Error reading from @%s.", u.GetNickName())
		}

		message = bytes.Trim(message, "\x00")

		input := string(message)
		//Remove new line
		if strings.Contains(input, "\n") == true {
			input = strings.TrimSpace(input)
		}

		handled, err := u.HandleNewInput(chatServer, input)
		if err != nil {
			logs.Errf(err, "Error reading input from user @%s.", u.nickName)
		}
		if handled {
			//If handled then continue reading
			continue
		}

		if input != "\n" && input != `` {
			u.SetIncoming(input)
		}
	}
}

// WriteTo to the user's connection and remembers the last message that was sent out
func (u *User) WriteTo() {
	for message := range u.outgoing {
		u.conn.Write([]byte(message + "\n"))
	}
}

// Disconnect a user from this server
func (u *User) Disconnect() error {
	logs.Infof("disconnecting=@%s", u.nickName)
	u.SetOutgoing("Good Bye, come back again.")

	// Wait 1 second before actually disconnecting
	<-time.After(time.Second * 1)
	return u.conn.Close()
}

// handleNewInput looks at user input and reacts to it
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
		u.SetOutgoing("You need to join a channel, use /join #channel or use /help for more info.")
		return false, errs.New("User is not in a channel, and input is not a command")
	}
	return u.handleBroadCastInput(chatServer, userInput)
}

func (u *User) handleBroadCastInput(chatServer *Server, userInput string) (bool, error) {
	channel, err := chatServer.GetChannel(u.GetChannel())
	if err != nil {
		return false, errs.Wrap(err, "Error getting channel from server")
	}
	channel.Broadcast(chatServer, `@`+u.nickName+`: `+userInput)
	return true, nil
}

func (u *User) handleCommandInput(chatServer *Server, input string) (bool, error) {
	userCommand, err := command.FromString(input)
	if err != nil {
		u.SetOutgoing("Invalid command, use /help for more info. Error:" + err.Error())
		logs.Errf(err, "Failed executing %s command by @s", input, u.nickName)
		return false, errs.Wrap(err, "Error getting command from user input.")
	}

	commandParams, err := u.GetCommandParams(chatServer, input, userCommand)
	if err != nil {
		return false, errs.Wrap(err, "Couldn't get command params")
	}

	if err = userCommand.Execute(*commandParams); err != nil {
		logs.Errf(err, "error \t @%s command=%s error=%s", u.nickName, input)
		return false, errs.Wrapf(err, "Couldn't execute command %q.", input)
	}

	logs.Infof("User @%s executed command: %s", u.nickName, input)
	return true, nil
}

func (u *User) GetCommandParams(chatServer *Server, input string, executable command.Executable) (*command.Params, error) {
	commandParams := &command.Params{
		User1:    u,
		RawInput: input,
		Server:   chatServer,
	}

	if executable.RequiresParam(`user2`) {
		nickname, err := command.ParseNickNameFomInput(input)
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
		channelName, err := command.ParseChannelFromInput(input)
		if err != nil {
			return nil, errs.Wrapf(err, "Could not find the required #channel in the input. Input %s", input)
		}
		channel, err := chatServer.GetChannel(channelName)
		if err != nil {
			return nil, errs.Wrapf(err, "Could not get channel from the server. Channel Name %s", channelName)
		}
		commandParams.Channel = channel
	}

	if executable.RequiresParam(`message`) {
		message, err := command.ParseMessageFromInput(input)
		if err != nil {
			return nil, errs.Wrap(err, "Could not required message in the input")
		}
		commandParams.Message = message
	}
	return commandParams, nil
}
