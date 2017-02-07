package chat

import (
	"sync"

	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/logs"
)

// Channel users can be in a channel and chat with each other
type Channel struct {
	name     string
	lockName *sync.Mutex

	users     map[string]bool
	lockUsers *sync.Mutex
}

// NewChannel returns a channel
func NewChannel() *Channel {
	return &Channel{
		users:     make(map[string]bool),
		lockName:  new(sync.Mutex),
		lockUsers: new(sync.Mutex),
	}
}

// AddUser adds a user to this channel
func (c *Channel) AddUser(nickName string) {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	c.users[nickName] = true
}

// RemoveUser removes a user from this server
func (c *Channel) RemoveUser(nickName string) {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	delete(c.users, nickName)
}

// GetName gets a channel's name
func (c *Channel) GetName() string {
	c.lockName.Lock()
	defer c.lockName.Unlock()
	return c.name
}

// SetName sets a channel's name
func (c *Channel) SetName(channelName string) {
	c.lockName.Lock()
	defer c.lockName.Unlock()
	c.name = channelName
}

// GetUserCount returns the number of connected users to this channel
func (c *Channel) GetUserCount() int {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	return len(c.users)
}

// GetUsers returns nicknames of connected users
func (c *Channel) GetUsers() map[string]bool {
	c.lockUsers.Lock()
	defer c.lockUsers.Unlock()
	return cloneNickNames(c.users)
}

// Broadcast sends a message to every user in a channel
func (c *Channel) Broadcast(chatServer *Server, message string) {
	users := c.GetUsers()
	for nickName := range users {
		user, err := chatServer.GetUser(nickName)
		// User may no longer be connected to the chat server
		if err != nil {
			c.RemoveUser(nickName)
			logs.ErrIfErrf(err, "User %s is in channel %s but not on connected to the server", nickName, c.GetName())
			continue
		}
		go user.SetOutgoing(plugins.UserOutPutTChannel, message)
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
