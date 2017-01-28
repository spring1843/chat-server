// Package chat implements a chat server
// It aims to handle connections, manage users and channels and allow execution of chat commands
package chat

import (
	"errors"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

type (
	// Service  keeps listening for connections, it contains users and channels
	Service struct {
		Connection chan Connection
		Logger     *log.Logger
		Channels   map[string]*Channel
		Users      map[string]*User
		Incoming   chan string
		Outgoing   chan string
		CanLog     bool
		lock       *sync.Mutex
	}
)

// NewService returns a new instance of the Server
func NewService() Server {
	runningServer := &Service{
		Connection: make(chan Connection),
		Channels:   make(map[string]*Channel),
		Users:      make(map[string]*User),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
		Logger:     new(log.Logger),
		CanLog:     false,
		lock:       new(sync.Mutex),
	}
	return runningServer
}

// Listen Makes this server start listening to connections, when a user is connected he or she is welcomed
func (s *Service) Listen() {
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
func (s *Service) SetLogFile(file io.Writer) {
	logger := new(log.Logger)
	logger.SetOutput(file)
	s.Logger = logger
	s.CanLog = true
}

// LogPrintf is a centralized logging function, so that all logs go to the same file and they all have time stamps
// Ads a time stamp to every log entry
// For readability start the message with a category followed by \t
func (s *Service) LogPrintf(format string, v ...interface{}) {
	if s.CanLog != true {
		return
	}
	now := time.Now()
	s.Logger.Printf(now.Format(time.UnixDate)+"\t"+format, v...)
}

// AddUser to this server
func (s *Service) AddUser(user *User) {
	s.lock.Lock()
	s.Users[user.nickName] = user
	s.lock.Unlock()
}

// RemoveUser from this server
func (s *Service) RemoveUser(nickName string) error {
	s.lock.Lock()
	delete(s.Users, nickName)
	s.lock.Unlock()
	return nil
}

// RemoveUserFromChannel removes a user from a channel
func (s *Service) RemoveUserFromChannel(nickName, channelName string) error {
	channel, err := s.GetChannel(channelName)
	if err != nil {
		return err
	}

	channel.RemoveUser(nickName)
	return nil
}

// GetUser gets a connected user
func (s *Service) GetUser(nickName string) (*User, error) {
	s.lock.Lock()
	if _, ok := s.Users[nickName]; ok {
		user := s.Users[nickName]
		s.lock.Unlock()
		return user, nil
	}
	s.lock.Unlock()
	return nil, errors.New(`User @` + nickName + ` not connected`)
}

// IsUserConnected checks to see if a user with the given nickname is connected to this server or not
func (s *Service) IsUserConnected(nickName string) bool {
	_, err := s.GetUser(nickName)
	if err != nil {
		return false
	}
	return true
}

// ReceiveConnection is used when there's a new connection
func (s *Service) ReceiveConnection(conn Connection) {
	s.Connection <- conn
}

// GetChannel gets a channel from the given channelName
func (s *Service) GetChannel(channelName string) (*Channel, error) {
	s.lock.Lock()
	if _, ok := s.Channels[channelName]; ok {
		channel := s.Channels[channelName]
		s.lock.Unlock()
		return channel, nil
	}
	s.lock.Unlock()

	return nil, errors.New(`Channel #` + channelName + ` does not exist on this server`)
}

// GetChannelCount returns the number of channels on this server
func (s *Service) GetChannelCount() int {
	s.lock.Lock()
	count := len(s.Channels)
	s.lock.Unlock()
	return count
}

// AddChannel adds a channel to this server
func (s *Service) AddChannel(channelName string) *Channel {
	channel := NewChannel()
	channel.Name = channelName

	s.lock.Lock()
	s.Channels[channelName] = channel
	s.lock.Unlock()

	return channel
}

// ConnectedUsersCount returns the number of connected users
func (s *Service) ConnectedUsersCount() int {
	s.lock.Lock()
	count := len(s.Users)
	s.lock.Unlock()
	return count
}

// Broadcast sends a message to every user connected to the server
func (s *Service) Broadcast(message string) {
	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	s.lock.Lock()
	users := s.Users
	s.lock.Unlock()

	for nickName := range users {
		user, err := s.GetUser(nickName)
		// User may no longer be connected to the chat server
		if err != nil {
			continue
		}
		user.outgoing <- message
	}
}

// WelcomeNewUser shows a welcome message to a new user and makes a new user entity by asking the new user to pick a nickname
func (s *Service) WelcomeNewUser(connection Connection) {
	s.LogPrintf("connection \t New connection from address=%s", connection.RemoteAddr().String())

	user := NewConnectedUser(s, connection)
	user.outgoing <- "Welcome to chat server. There are " + strconv.Itoa(s.ConnectedUsersCount()) + " other users on this server. please enter a nickname"

	nickName := <-user.incoming

	for s.IsUserConnected(nickName) {
		user.outgoing <- "Another user with this nickname is connected to this server, Please enter a different nickname"
		nickName = <-user.incoming
	}

	user.nickName = nickName
	s.AddUser(user)
	s.LogPrintf("connection \t address=%s authenticated=@%s", connection.RemoteAddr().String(), nickName)

	user.outgoing <- "Thanks " + user.nickName + ", now please type /join #channel to join a channel or /help to get all commands"
}
