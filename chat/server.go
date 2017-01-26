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

var RunningServer = NewService()

type (
	// Service  keeps listening for connections, it contains users and channels
	Service struct {
		Connection chan Connection
		Logger     *log.Logger
		Channels   []*Channel
		Users      map[string]*User
		Incoming   chan string
		Outgoing   chan string
		CanLog     bool
		lock       *sync.Mutex
	}
)

// NewService returns a new instance of the Server
func NewService() *Service {
	return &Service{
		Connection: make(chan Connection),
		Users:      make(map[string]*User),
		Channels:   make([]*Channel, 0),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
		Logger:     new(log.Logger),
		CanLog:     false,
		lock:       new(sync.Mutex),
	}
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
	s.Users[user.NickName] = user
	s.lock.Unlock()
}

// RemoveUser from this server
func (s *Service) RemoveUser(nickName string) error {
	s.lock.Lock()
	if _, ok := s.Users[nickName]; !ok {
		return errors.New("Can not remove user, nickname " + nickName + " does not exist")
	}
	delete(s.Users, nickName)
	s.lock.Unlock()
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

// GetChannel gets a channel from the given channelName
func (s *Service) GetChannel(channelName string) (*Channel, error) {
	for _, channel := range s.Channels {
		if channel.Name == channelName {
			return channel, nil
		}
	}
	return nil, errors.New(`Channel #` + channelName + ` does not exist on this server`)
}

// AddChannel adds a channel to this server
func (s *Service) AddChannel(channelName string) *Channel {
	channel := NewChannel()
	channel.Name = channelName
	s.Channels = append(s.Channels, channel)
	return channel
}

func (s *Service) ConnectedUsersCount() int {
	s.lock.Lock()
	count := len(s.Users)
	s.lock.Unlock()
	return count
}

// WelcomeNewUser shows a welcome message to a new user and makes a new user entity by asking the new user to pick a nickname
func (s *Service) WelcomeNewUser(connection Connection) {
	s.LogPrintf("connection \t New connection from address=%s", connection.RemoteAddr().String())

	user := NewUser(connection)
	user.Outgoing <- "Welcome to chat server. There are " + strconv.Itoa(s.ConnectedUsersCount()) + " other users on this server. please enter a nickname"

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
