package chat

import (
	"bytes"
	"errors"
	"strings"
)

// User is temporarily in connected to a chat server, and can be in certain channels
type User struct {
	Connection          Connection
	NickName            string
	Channel             *Channel
	IgnoreList          map[string]bool
	Incoming            chan string
	Outgoing            chan string
	LastOutGoingMessage string
	LastIncomingMessage string
}

// NewUser returns a new User
func NewUser(connection Connection) *User {
	User := &User{
		Connection: connection,
		NickName:   ``,
		Channel:    nil,
		IgnoreList: make(map[string]bool),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
	}
	User.Listen()

	return User
}

// Write to the user's connection and remembers the last message that was sent out
func (u *User) Write() {
	for message := range u.Outgoing {
		u.Connection.Write([]byte(message + "\n"))
		u.LastOutGoingMessage = message
	}
}

// Read and interprets a message from a user
func (u *User) Read() {
	for {
		message := make([]byte, 256)
		u.Connection.Read(message)

		message = bytes.Trim(message, "\x00")

		input := string(message)
		//Remove new line
		if strings.Contains(input, "\n") == true {
			input = strings.TrimSpace(input)
		}

		if u.handleNewInput(input) {
			//If handled then continue reading
			continue
		}

		if input != "\n" && input != `` {
			u.Incoming <- input
		}
	}
}

// Listen starts reading from and writing to a user
func (u *User) Listen() {
	go u.Read()
	go u.Write()
}

// Ignore a user
func (u *User) Ignore(nickName string) {
	u.IgnoreList[nickName] = true
}

// HasIgnored checks to see if a user has ignored another user or not
func (u *User) HasIgnored(nickName string) bool {
	if _, ok := u.IgnoreList[nickName]; ok {
		return true
	}
	return false
}

// Disconnect a user from this server
func (u *User) Disconnect() {
	RunningServer.LogPrintf("connection \t disconnecting=@%s", u.NickName)
	u.Connection.Close()
}

// Checks to see if a new input from user is a command
// If it is a command then it tries executing func
// If it's not a command then it will output to the channel
func (u *User) handleNewInput(input string) bool {
	if command, err := GetCommand(input); err == nil && command != nil {
		err = u.ExecuteCommand(input, command)
		if err != nil {
			u.Outgoing <- `Could not execute command. Error:` + err.Error()
			RunningServer.LogPrintf("error \t failed @%s's command %s", u.NickName, input)
		}
		return true
	}

	if u.Channel != nil {
		RunningServer.LogPrintf("message \t @%s in #%s message=%s", u.NickName, u.Channel.Name, input)
		u.Channel.Broadcast(RunningServer, `@`+u.NickName+`: `+input)
		return true
	}

	return false
}

// ExecuteCommand Executes a given command
// First it finds all the required parameters from the input and populates them
func (u *User) ExecuteCommand(input string, command Executable) error {
	commandParams := CommandParams{
		user1:    u,
		rawInput: input,
		server:   RunningServer,
	}

	chatCommand := command.getChatCommand()

	if chatCommand.DoesCommandRequireParam(`user2`) == true {
		nickname, err := chatCommand.ParseNickNameFomInput(input)
		if err != nil {
			return errors.New("Could not find the required @nickname in the input")
		}

		user2, err := RunningServer.GetUser(nickname)
		if err != nil {
			return err
		}
		commandParams.user2 = user2
	}

	if chatCommand.DoesCommandRequireParam(`channel`) == true {
		channelName, err := chatCommand.ParseChannelFromInput(input)
		if err != nil {
			return errors.New("Could not find the required #channel in the input")
		}
		channel, err := RunningServer.GetChannel(channelName)
		if err == nil {
			commandParams.channel = channel
		}
	}

	if chatCommand.DoesCommandRequireParam(`message`) == true {
		message, err := chatCommand.ParseMessageFromInput(input)
		if err != nil {
			return errors.New("Could not required message in the input")
		}
		commandParams.message = message
	}

	err := command.Execute(commandParams)
	if err != nil {
		RunningServer.LogPrintf("error \t @%s command=%s error=%s", u.NickName, input, err.Error())
		return err
	}
	RunningServer.LogPrintf("command \t @%s command=%s", u.NickName, input)

	return nil
}
