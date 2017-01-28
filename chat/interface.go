package chat

import (
	"io"
	"net"
)

type (
	// Server is an interface for a chat server
	Server interface {
		Listen()
		ReceiveConnection(conn Connection)
		Broadcast(message string)

		SetLogFile(file io.Writer)
		LogPrintf(format string, v ...interface{})

		AddUser(user *User)
		RemoveUser(nickName string) error
		GetUser(nickName string) (*User, error)
		IsUserConnected(string) bool
		ConnectedUsersCount() int
		WelcomeNewUser(connection Connection)

		GetChannel(string) (*Channel, error)
		AddChannel(channelName string) *Channel
		GetChannelCount() int
		RemoveUserFromChannel(nickName, channelName string) error
	}
	// Chan is an interface for a chat channel
	Chan interface {
		AddUser(nickName string)
		RemoveUser(nickName string)
		Broadcast(chatServer Server, message string)
		GetName() string
		GetUserCount() int
	}
	// Chatter is an interface for a chat user
	Chatter interface {
		SetOutgoing(message string)
		GetChannel() string
		GetNickName() string
		Ignore(nickName string)
		SetChannel(name string)
		HasIgnored(nickName string) bool
		Disconnect(chatServer Server) error
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
	// Connection is an interface for a network connection
	Connection interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
		Close() error
		RemoteAddr() net.Addr
	}
)
