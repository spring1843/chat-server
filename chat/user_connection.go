package chat

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spring1843/chat-server/plugins/command"
	"github.com/spring1843/chat-server/plugins/errs"
	"github.com/spring1843/chat-server/plugins/logs"
)

// NewConnectedUser returns a new User with a connection
func NewConnectedUser(chatServer *Server, connection Connection) *User {
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
	fmt.Printf("Getting incomming\n")
	return <-u.incoming
}

// SetIncoming sets an incoming message from the user
func (u *User) SetIncoming(message string) {
	fmt.Printf("Setting incomming\n")
	u.incoming <- message
}

// ReadFrom reads data from users and lets chat server interpret it
func (u *User) ReadFrom(chatServer *Server) {
	for {
		message := make([]byte, 256)
		if _, err := u.conn.Read(message); err != nil {
			logs.Errf(err, "Error reading from @%s.", u.GetNickName())
		}

		message = bytes.Trim(message, "\x00")

		input := string(message)
		//Remove new line
		if strings.Contains(input, "\n") == true {
			input = strings.TrimSpace(input)
		}

		handled, err := u.handleNewInput(chatServer, input)
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

// Write to the user's connection and remembers the last message that was sent out
func (u *User) WriteTo() {
	for message := range u.outgoing {
		u.conn.Write([]byte(message + "\n"))
	}
}

// Disconnect a user from this server
func (u *User) Disconnect(chatServer *Server) error {
	logs.Infof("connection \t disconnecting=@%s", u.nickName)
	return u.conn.Close()
}

// Checks to see if a new input from user is a command
// If it is a command then it tries executing func
// If it's not a command then it will output to the channel
func (u *User) handleNewInput(chatServer *Server, input string) (bool, error) {
	if command, err := command.GetCommand(input); err == nil && command != nil {
		err = u.ExecuteCommand(chatServer, input, command)
		if err != nil {
			u.outgoing <- `Could not execute command. Error:` + err.Error()
			err = errs.Wrap(err, "error handling input")
			logs.Errf(err, "Failed executing %s command by @s", input, u.nickName)
			return false, err
		}
		return true, nil
	}

	if u.GetChannel() != "" {
		logs.Infof("message \t @%s in #%s message=%s", u.nickName, u.GetChannel(), input)
		channel, err := chatServer.GetChannel(u.GetChannel())
		if err != nil {
			return false, errs.Wrap(err, "Error getting channel from server")
		}
		channel.Broadcast(chatServer, `@`+u.nickName+`: `+input)
		return true, nil
	}

	return false, nil
}

// ExecuteCommand Executes a given command
// First it finds all the required parameters from the input and populates them
func (u *User) ExecuteCommand(chatServer *Server, input string, executable command.Executable) error {
	commandParams := command.CommandParams{
		User1:    u,
		RawInput: input,
		Server:   chatServer,
	}

	chatCommand := executable.GetChatCommand()

	if chatCommand.RequiresParam(`user2`) {
		nickname, err := chatCommand.ParseNickNameFomInput(input)
		if err != nil {
			return errs.Wrap(err, "Could not find the required @nickname in the input")
		}

		user2, err := chatServer.GetUser(nickname)
		if err != nil {
			return errs.Wrap(err, "User "+nickname+" + is not connected to this server")
		}
		commandParams.User2 = user2
	}

	if chatCommand.RequiresParam(`channel`) {
		channelName, err := chatCommand.ParseChannelFromInput(input)
		if err != nil {
			return errs.Wrap(err, "Could not find the required #channel in the input")
		}
		channel, err := chatServer.GetChannel(channelName)
		if err == nil {
			commandParams.Channel = channel
		}
	}

	if chatCommand.RequiresParam(`message`) {
		message, err := chatCommand.ParseMessageFromInput(input)
		if err != nil {
			return errs.Wrap(err, "Could not required message in the input")
		}
		commandParams.Message = message
	}

	err := executable.Execute(commandParams)
	if err != nil {
		logs.Errf(err, "error \t @%s command=%s error=%s", u.nickName, input)
		return err
	}

	logs.Infof("command \t @%s command=%s", u.nickName, input)

	return nil
}
