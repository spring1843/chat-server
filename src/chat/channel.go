package chat

import (
	"sync"

	"github.com/spring1843/chat-server/src/shared/logs"
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
	return cloneNickNames(c.Users)
}

// Broadcast sends a message to every user in a channel
func (c *Channel) Broadcast(chatServer *Server, message string) {
	users := c.GetUsers()
	for nickName := range users {
		user, err := chatServer.GetUser(nickName)
		// User may no longer be connected to the chat server
		if err != nil {
			c.RemoveUser(nickName)
			logs.Errf(err, "User %s is in channel %s but not on connected to the server", user.GetNickName(), c.GetName())
			continue
		}
		user.SetOutgoing(message)
	}
}

// getNickNames turns a map of nicknames to slice of strings
// This copy is so that we can unlock faster and avoid race
// Since we are looping again through nickNames next and making another blocking call
func cloneNickNames(users map[string]bool) map[string]bool {
	nickNames := make(map[string]bool)
	for nickName := range users {
		nickNames[nickName] = true
	}
	return nickNames
}
