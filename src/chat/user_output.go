package chat

import (
	"fmt"
	"time"
)

const outMessagePattern = "%d %02d %s"

// SetOutgoing sets an outgoing message to the user
func (u *User) SetOutgoing(messageType int, message string) {
	u.outgoing <- fmt.Sprintf(outMessagePattern, time.Now().Unix(), messageType, message)
}

// SetOutgoingf sets an outgoing message to the user
func (u *User) SetOutgoingf(messageType int, format string, a ...interface{}) {
	u.SetOutgoing(messageType, fmt.Sprintf(format, a...))
}

// GetOutgoing gets the outgoing message for a user
func (u *User) GetOutgoing() string {
	return <-u.outgoing
}
