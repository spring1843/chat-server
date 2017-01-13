// Package chat implements a chat server
// It aims to handle connections, manage users and channels and allow execution of chat commands
package chat

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type (
	// Server  keeps listening for connections, it contains users and channels
	Server struct {
		Connection chan Connection
		Logger     *log.Logger
		Channels   []*Channel
		Users      []*User
		Incoming   chan string
		Outgoing   chan string
		CanLog     bool
	}
	// Connection defines behaviors of a connected user
	Connection interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
		Close() error
		RemoteAddr() net.Addr
	}
)

// NewServer returns a new instance of the Server
func NewServer() *Server {
	return &Server{
		Connection: make(chan Connection),
		Users:      make([]*User, 0),
		Channels:   make([]*Channel, 0),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
		Logger:     new(log.Logger),
		CanLog:     false,
	}
}

// Listen Makes this server start listening to connections, when a user is connected he or she is welcomed
func (s *Server) Listen() {
	go func() {
		for {
			for connection := range s.Connection {
				go s.WelcomeNewUser(connection)
			}
		}
	}()
}

// SetLogFile a log file for this server and makes this server able to log
// Use server.Log() to send logs to this file
func (s *Server) SetLogFile(file io.Writer) {
	logger := new(log.Logger)
	logger.SetOutput(file)
	s.Logger = logger
	s.CanLog = true
}

// LogPrintf is a centralized logging function, so that all logs go to the same file and they all have time stamps
// Ads a time stamp to every log entry
// For readability start the message with a category followed by \t
func (s *Server) LogPrintf(format string, v ...interface{}) {
	if s.CanLog != true {
		return
	}
	now := time.Now()
	s.Logger.Printf(now.Format(time.UnixDate)+"\t"+format, v...)
}

// AddUser to this server
func (s *Server) AddUser(user *User) {
	user.SetServer(s)
	s.Users = append(s.Users, user)
}

// RemoveUser from this server
func (s *Server) RemoveUser(user *User) error {
	i := -1
	for _, user := range s.Users {
		i++
		if user.NickName == user.NickName {
			break
		}
	}
	if i == -1 {
		return errors.New(`Did not find user to remove`)
	}
	copyUsers := s.Users
	copyUsers = append(copyUsers[:i], copyUsers[i+1:]...)
	s.Users = copyUsers

	return nil
}

// GetUser gets a connected user
func (s *Server) GetUser(nickName string) (*User, error) {
	for _, user := range s.Users {
		if user.NickName == nickName {
			return user, nil
		}
	}
	return nil, errors.New(`User @` + nickName + ` not connected`)
}

// IsUserConnected checks to see if a user with the given nickname is connected to this server or not
func (s *Server) IsUserConnected(nickName string) bool {
	_, err := s.GetUser(nickName)
	if err != nil {
		return false
	}
	return true
}

// GetChannel gets a channel from the given channelName
func (s *Server) GetChannel(channelName string) (*Channel, error) {
	for _, channel := range s.Channels {
		if channel.Name == channelName {
			return channel, nil
		}
	}
	return nil, errors.New(`Channel #` + channelName + ` does not exist on this server`)
}

// AddChannel adds a channel to this server
func (s *Server) AddChannel(channelName string) *Channel {
	channel := NewChannel()
	channel.Name = channelName
	s.Channels = append(s.Channels, channel)
	return channel
}

// WelcomeNewUser shows a welcome message to a new user and makes a new user entity by asking the new user to pick a nickname
func (s *Server) WelcomeNewUser(connection Connection) {
	s.LogPrintf("connection \t New connection from address=%s", connection.RemoteAddr().String())

	user := NewUser(connection)
	user.Outgoing <- "Welcome to chat server. There are " + strconv.Itoa(len(s.Users)) + " other users on this server. please enter a nickname"

	nickName := <-user.Incoming

	for s.IsUserConnected(nickName) {
		user.Outgoing <- "Another user with this nickname is connected to this server, Please enter a different nickname"
		nickName = <-user.Incoming
	}

	user.NickName = nickName
	s.AddUser(user)
	s.LogPrintf("connection \t address=%s authenticated=@%s", connection.RemoteAddr().String(), nickName)

	user.Outgoing <- "Thanks " + user.NickName + ", now please type /join #channel to join a channel or /help to get all commands"
}
