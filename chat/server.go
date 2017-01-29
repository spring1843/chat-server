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

// Server  keeps listening for connections, it contains users and channels
type Server struct {
	Connection chan Connection
	Logger     *log.Logger
	Channels   map[string]*Channel
	Users      map[string]*User
	Incoming   chan string
	Outgoing   chan string
	CanLog     bool
	lock       *sync.Mutex
}

// NewServer returns a new instance of the chat server
func NewServer() *Server {
	server := &Server{
		Connection: make(chan Connection),
		Channels:   make(map[string]*Channel),
		Users:      make(map[string]*User),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
		Logger:     new(log.Logger),
		CanLog:     false,
		lock:       new(sync.Mutex),
	}
	return server
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
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Users[user.nickName] = user
}

// RemoveUser from this server
func (s *Server) RemoveUser(nickName string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.Users[nickName]; !ok {
		return errors.New("User " + nickName + " is not connected to this server")
	}
	delete(s.Users, nickName)
	return nil
}

// RemoveUserFromChannel removes a user from a channel
func (s *Server) RemoveUserFromChannel(nickName, channelName string) error {
	channel, err := s.GetChannel(channelName)
	if err != nil {
		return err
	}

	channel.RemoveUser(nickName)
	return nil
}

// GetUser gets a connected user
func (s *Server) GetUser(nickName string) (*User, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.Users[nickName]; ok {
		user := s.Users[nickName]
		return user, nil
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

// ReceiveConnection is used when there's a new connection
func (s *Server) ReceiveConnection(conn Connection) {
	s.Connection <- conn
}

// GetChannel gets a channel from the given channelName
func (s *Server) GetChannel(channelName string) (*Channel, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.Channels[channelName]; ok {
		channel := s.Channels[channelName]
		return channel, nil
	}

	return nil, errors.New(`Channel #` + channelName + ` does not exist on this server`)
}

// GetChannelCount returns the number of channels on this server
func (s *Server) GetChannelCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.Channels)
}

// AddChannel adds a channel to this server
func (s *Server) AddChannel(channelName string) {
	channel := NewChannel()
	channel.Name = channelName

	s.lock.Lock()
	defer s.lock.Unlock()
	s.Channels[channelName] = channel
}

// ConnectedUsersCount returns the number of connected users
func (s *Server) ConnectedUsersCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.Users)
}

// Broadcast sends a message to every user connected to the server
func (s *Server) Broadcast(message string) {
	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	s.lock.Lock()
	defer s.lock.Unlock()
	users := s.Users

	for nickName := range users {
		user, err := s.GetUser(nickName)
		// User may no longer be connected to the chat server
		if err != nil {
			continue
		}
		user.outgoing <- message
	}
}

// DisconnectUser disconnects a user from this server
func (s *Server) DisconnectUser(nickName string) error {
	user, err := s.GetUser(nickName)
	if err != nil {
		return err
	}
	return user.Disconnect(s)
}

func (s *Server) BroadcastInChannel(channelName, message string) error {
	channel, err := s.GetChannel(channelName)
	if err != nil {
		return err
	}

	channel.Broadcast(s, message)
	return nil
}

// WelcomeNewUser shows a welcome message to a new user and makes a new user entity by asking the new user to pick a nickname
func (s *Server) WelcomeNewUser(connection Connection) {
	s.LogPrintf("connection \t New connection from address=%s", connection.RemoteAddr().String())

	user := NewConnectedUser(s, connection)
	user.SetOutgoing("Welcome to chat server. There are " + strconv.Itoa(s.ConnectedUsersCount()) + " other users on this server. please enter a nickname")

	nickName := user.GetIncoming()

	for s.IsUserConnected(nickName) {
		user.SetOutgoing("Another user with this nickname is connected to this server, Please enter a different nickname")
		nickName = user.GetIncoming()
	}

	user.SetNickName(nickName)
	s.AddUser(user)
	s.LogPrintf("connection \t address=%s authenticated=@%s", connection.RemoteAddr().String(), nickName)

	user.SetOutgoing("Thanks " + user.nickName + ", now please type /join #channel to join a channel or /help to get all commands")
}
