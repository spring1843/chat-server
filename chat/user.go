package chat

import (
	"bytes"
	"errors"
	"strings"
)

// A user entity
type User struct {
	Connection          ChatConnection
	NickName            string
	Server              *Server
	Channel             *Channel
	IgnoreList          []User
	Incoming            chan string
	Outgoing            chan string
	LastOutGoingMessage string
	LastIncomingMessage string
}

func NewUser(connection ChatConnection) *User {
	User := &User{
		Connection: connection,
		NickName:   ``,
		Server:     nil,
		Channel:    nil,
		IgnoreList: make([]User, 0),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
	}

	User.Listen()

	return User
}

// Writes to the user's connection and remembers the last message that was sent out
func (u *User) Write() {
	for message := range u.Outgoing {
		u.Connection.Write([]byte(message + "\n"))
		u.LastOutGoingMessage = message
	}
}

// Reads and interprets a message from a user
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

		if u.handleNewInput(input) == true {
			//If handled then continue reading
			continue
		}

		if input != "\n" && input != `` {
			u.Incoming <- input
		}

	}
}

// Sets the server this user belongs to
func (u *User) SetServer(server *Server) {
	u.Server = server
}

// Starts reading from and writing to a user
func (u *User) Listen() {
	go u.Read()
	go u.Write()
}

// Ignores a user
func (u *User) Ignore(user User) {
	u.IgnoreList = append(u.IgnoreList, user)
}

// Checks to see if a user has ignored another user or not
func (u *User) HasIgnored(user User) bool {
	for _, ignoredUser := range u.IgnoreList {
		if ignoredUser.NickName == user.NickName {
			return true
		}
	}
	return false
}

// Disconnects a user from this server
func (u *User) Disconnect() {
	u.Server.LogPrintf("connection \t disconnecting=@%s", u.NickName)
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
			u.Server.LogPrintf("error \t failed @%s's command %s", u.NickName, input)
		}
		return true
	}

	if u.Channel != nil {
		u.Server.LogPrintf("message \t @%s in #%s message=%s", u.NickName, u.Channel.Name, input)
		u.Channel.Broadcast(`@` + u.NickName + `: ` + input)
		return true
	}

	return false
}

// Executes a given command
// First it finds all the required parameters from the input and populates them
func (u *User) ExecuteCommand(input string, command Executable) error {
	commandParams := CommandParams{
		user1:    u,
		rawInput: input,
	}

	chatCommand := command.getChatCommand()

	if chatCommand.DoesCommandRequireParam(`user2`) == true {
		nickname, err := chatCommand.ParseNickNameFomInput(input)
		if err != nil {
			return errors.New("Could not find the required @nickname in the input")
		}

		user2, err := u.Server.GetUser(nickname)
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
		channel, err := u.Server.GetChannel(channelName)
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
		u.Server.LogPrintf("error \t @%s command=%s error=%s", u.NickName, input, err.Error())
		return err
	}
	u.Server.LogPrintf("command \t @%s command=%s", u.NickName, input)

	return nil
}
