// Package chat implements a chat server
// It aims to handle connections, manage users and channels and allow execution of chat commands
package chat

import (
	"errors"
	"sync"
	"time"

	"github.com/spring1843/chat-server/drivers"
)

// Server  keeps listening for connections, it contains users and channels
type Server struct {
	Connection chan drivers.Connection
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
		Connection: make(chan drivers.Connection),
		Channels:   make(map[string]*Channel),
		Users:      make(map[string]*User),
		Incoming:   make(chan string),
		Outgoing:   make(chan string),
		CanLog:     false,
		lock:       new(sync.Mutex),
	}
	return server
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

// ConnectedUsersCount returns the number of connected users
func (s *Server) ConnectedUsersCount() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.Users)
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

// BroadcastInChannel broadcasts a message to all the users in a channel
func (s *Server) BroadcastInChannel(channelName, message string) error {
	channel, err := s.GetChannel(channelName)
	if err != nil {
		return err
	}

	channel.Broadcast(s, message)
	return nil
}
