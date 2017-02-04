package chat

import (
	"fmt"
	"time"
)

const outMessagePattern = "%d %s %s"

// SetOutgoing sets an outgoing message to the user
func (u *User) SetOutgoing(messageType string, message string) {
	u.outgoing <- fmt.Sprintf(outMessagePattern, time.Now().Unix(), messageType, message)
}

// SetOutgoingf sets an outgoing message to the user
func (u *User) SetOutgoingf(messageType string, format string, a ...interface{}) {
	u.SetOutgoing(messageType, fmt.Sprintf(format, a...))
}

// GetOutgoing gets the outgoing message for a user
func (u *User) GetOutgoing() string {
	return <-u.outgoing
}
