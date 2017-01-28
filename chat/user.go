package chat

import (
	"bytes"
	"errors"
	"strings"
	"sync"
)

// User is temporarily in connected to a chat server, and can be in certain channels
type User struct {
	Connection          Connection
	NickName            string
	Channel             string
	IgnoreList          map[string]bool
	incoming            chan string
	outgoing            chan string
	LastOutGoingMessage string
	LastIncomingMessage string
	lock                *sync.Mutex
}

// NewUser returns a new new User
func NewUser(nickName string) *User {
	return &User{
		NickName:   nickName,
		Channel:    "",
		IgnoreList: make(map[string]bool),
		incoming:   make(chan string),
		outgoing:   make(chan string),
		lock:       new(sync.Mutex),
	}
}

// NewConnectedUser returns a new User with a connection
func NewConnectedUser(chatServer Server, connection Connection) *User {
	user := NewUser("")
	user.Connection = connection
	user.Listen(chatServer)
	return user
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

// GetChannel gets the current channel name for the user
func (u *User) GetChannel() string {
	u.lock.Lock()
	name := u.Channel
	u.lock.Unlock()
	return name
}

// SetChannel sets the current channel name for the user
func (u *User) SetChannel(name string) {
	u.lock.Lock()
	u.Channel = name
	u.lock.Unlock()
}

// Write to the user's connection and remembers the last message that was sent out
func (u *User) Write() {
	for message := range u.outgoing {
		u.Connection.Write([]byte(message + "\n"))
		u.LastOutGoingMessage = message
	}
}

// Read and interprets a message from a user
func (u *User) Read(chatServer Server) {
	for {
		message := make([]byte, 256)
		u.Connection.Read(message)

		message = bytes.Trim(message, "\x00")

		input := string(message)
		//Remove new line
		if strings.Contains(input, "\n") == true {
			input = strings.TrimSpace(input)
		}

		err, handled := u.handleNewInput(chatServer, input)
		if err != nil {
			chatServer.LogPrintf("Error reading input from user @%s. Error %s", u.NickName, err)
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

// Listen starts reading from and writing to a user
func (u *User) Listen(chatServer Server) {
	go u.Read(chatServer)
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
func (u *User) Disconnect(chatServer Server) error {
	chatServer.LogPrintf("connection \t disconnecting=@%s", u.NickName)
	return u.Connection.Close()
}

// Checks to see if a new input from user is a command
// If it is a command then it tries executing func
// If it's not a command then it will output to the channel
func (u *User) handleNewInput(chatServer Server, input string) (error, bool) {
	if command, err := GetCommand(input); err == nil && command != nil {
		err = u.ExecuteCommand(chatServer, input, command)
		if err != nil {
			u.outgoing <- `Could not execute command. Error:` + err.Error()
			chatServer.LogPrintf("error \t failed @%s's command %s", u.NickName, input)
		}
		return nil, true
	}

	if u.GetChannel() != "" {
		chatServer.LogPrintf("message \t @%s in #%s message=%s", u.NickName, u.GetChannel(), input)
		channel, err := chatServer.GetChannel(u.GetChannel())
		if err != nil {
			return err, false
		}
		channel.Broadcast(chatServer, `@`+u.NickName+`: `+input)
		return nil, true
	}

	return nil, false
}

// ExecuteCommand Executes a given command
// First it finds all the required parameters from the input and populates them
func (u *User) ExecuteCommand(chatServer Server, input string, command Executable) error {
	commandParams := CommandParams{
		user1:    u,
		rawInput: input,
		server:   chatServer,
	}

	chatCommand := command.getChatCommand()

	if chatCommand.RequiresParam(`user2`) {
		nickname, err := chatCommand.ParseNickNameFomInput(input)
		if err != nil {
			return errors.New("Could not find the required @nickname in the input")
		}

		user2, err := chatServer.GetUser(nickname)
		if err != nil {
			return errors.New("User " + nickname + " + is not connected to this server")
		}
		commandParams.user2 = user2
	}

	if chatCommand.RequiresParam(`channel`) {
		channelName, err := chatCommand.ParseChannelFromInput(input)
		if err != nil {
			return errors.New("Could not find the required #channel in the input")
		}
		channel, err := chatServer.GetChannel(channelName)
		if err == nil {
			commandParams.channel = channel
		}
	}

	if chatCommand.RequiresParam(`message`) {
		message, err := chatCommand.ParseMessageFromInput(input)
		if err != nil {
			return errors.New("Could not required message in the input")
		}
		commandParams.message = message
	}

	err := command.Execute(commandParams)
	if err != nil {
		chatServer.LogPrintf("error \t @%s command=%s error=%s", u.NickName, input, err.Error())
		return err
	}

	chatServer.LogPrintf("command \t @%s command=%s", u.NickName, input)

	return nil
}
