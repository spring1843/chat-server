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
	c.Users[nickName] = true
	c.lock.Unlock()
}

// RemoveUser removes a user from this server
func (c *Channel) RemoveUser(nickName string) {
	c.lock.Lock()
	delete(c.Users, nickName)
	c.lock.Unlock()
}

// Broadcast sends a message to every user in a channel
func (c *Channel) Broadcast(chatServer Server, message string) {
	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	c.lock.Lock()
	users := c.Users
	c.lock.Unlock()

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
