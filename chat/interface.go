package chat

import (
	"io"
	"net"
)

type (
	// Server is an interface for the chat server
	Server interface {
		Listen()
		SetLogFile(file io.Writer)
		LogPrintf(format string, v ...interface{})
		AddUser(user *User)
		RemoveUser(*User) error
		GetUser(string) (*User, error)
		IsUserConnected(string) bool
		GetChannel(string) (*Channel, error)
		AddChannel(channelName string) *Channel
		WelcomeNewUser(connection Connection)
	}
	// Connection is an interface for a network connection
	Connection interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
		Close() error
		RemoteAddr() net.Addr
	}
	// Chatter is an interface for users
	Chatter interface {
		Write()
		Read()
		SetServer(server *Service)
		Listen()
		Ignore(nickName string)
		HasIgnored(nickName string) bool
		Disconnect()
		ExecuteCommand(input string, command Executable)
	}
	// Executable is an interface for chat commands
	Executable interface {
		Execute(params CommandParams) error
		getChatCommand() Command
		ParseNickNameFomInput(input string) (string, error)
		ParseChannelFromInput(input string) (string, error)
		ParseMessageFromInput(input string) (string, error)
		ParseCommandFromInput(input string) (string, error)
	}
)
