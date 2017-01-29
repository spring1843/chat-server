package chat

import (
	"sync"
	"time"
)

// Channel users can be in a channel and chat with each other
type Channel struct {
	Name     string
	lockName *sync.Mutex

	Users     map[string]bool
	lockUsers *sync.Mutex
}

// NewChannel returns a channel
func NewChannel() *Channel {
	return &Channel{
		Users:     make(map[string]bool),
		lockName:  new(sync.Mutex),
		lockUsers: new(sync.Mutex),
	}
}

// AddUser adds a user to this channel
func (c *Channel) AddUser(nickName string) {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	c.Users[nickName] = true
}

// RemoveUser removes a user from this server
func (c *Channel) RemoveUser(nickName string) {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	delete(c.Users, nickName)
}

// GetName gets a channel's name
func (c *Channel) GetName() string {
	c.lockName.Lock()
	defer c.lockName.Unlock()
	return c.Name
}

// SetName sets a channel's name
func (c *Channel) SetName(channelName string) {
	c.lockName.Lock()
	defer c.lockName.Unlock()
	c.Name = channelName
}

// GetUserCount returns the number of connected users to this channel
func (c *Channel) GetUserCount() int {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	return len(c.Users)
}

// GetUsers returns nicknames of connected users
func (c *Channel) GetUsers() map[string]bool {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	return c.Users
}

// Broadcast sends a message to every user in a channel
func (c *Channel) Broadcast(chatServer *Server, message string) {
	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	c.lockUsers.Lock()
	users := c.Users
	c.lockUsers.Unlock()

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
