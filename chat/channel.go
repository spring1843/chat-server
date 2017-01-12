package chat

import (
	"errors"
	"time"
)

// Channel users can be in a channel and chat with each other
type Channel struct {
	Name  string
	Users []*User
}

// NewChannel returns a channel
func NewChannel() *Channel {
	channel := &Channel{
		Users: make([]*User, 0),
	}

	return channel
}

// AddUser adds a user to this channel
func (c *Channel) AddUser(user *User) {
	c.Users = append(c.Users, user)
}

// RemoveUser removes a user from this server
func (c *Channel) RemoveUser(user *User) error {
	i := -1
	for _, user := range c.Users {
		i++
		if user.NickName == user.NickName {
			break
		}
	}
	if i == -1 {
		return errors.New(`Did not find user to remove`)
	}
	copyUsers := c.Users
	copyUsers = append(copyUsers[:i], copyUsers[i+1:]...)
	c.Users = copyUsers

	return nil
}

// Broadcast sends a message to every user in the chat room
func (c *Channel) Broadcast(message string) {

	now := time.Now()
	message = now.Format(time.Kitchen) + `-` + message

	for _, user := range c.Users {
		user.Outgoing <- message
	}
}
