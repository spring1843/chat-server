package chat

import (
	"sync"
	"time"
)

// Channel users can be in a channel and chat with each other
type Channel struct {
	Name  string
	Users map[string]bool
	lock  sync.Mutex
}

// NewChannel returns a channel
func NewChannel() *Channel {
	return &Channel{
		Users: make(map[string]bool),
	}
}

// AddUser adds a user to this channel
func (c *Channel) AddUser(nickName string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Users[nickName] = true
}

// RemoveUser removes a user from this server
func (c *Channel) RemoveUser(nickName string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.Users, nickName)
}

// GetName gets a channel's name
func (c *Channel) GetName() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Name
}

// GetUserCount returns the number of connected users to this channel
func (c *Channel) GetUserCount() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.Users)
}

// GetUsers returns nicknames of connected users
func (c *Channel) GetUsers() map[string]bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Users
}

// Broadcast sends a message to every user in a channel
func (c *Channel) Broadcast(chatServer *Server, message string) {
	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	c.lock.Lock()
	defer c.lock.Unlock()
	users := c.Users

	for nickName := range users {
		user, err := chatServer.GetUser(nickName)
		// User may no longer be connected to the chat server
		if err != nil {
			c.RemoveUser(nickName)
			continue
		}
		user.outgoing <- message
	}
}
