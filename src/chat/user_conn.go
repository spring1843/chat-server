package chat

import (
	"bytes"
	"strings"
	"time"

	"github.com/spring1843/chat-server/src/drivers"
	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/logs"
)

// ReadConnectionLimitBytes is the maximum size of input we accept from user
// This is important to defend against DOS attacks
var ReadConnectionLimitBytes = 100000 // 100KB

// NewConnectedUser returns a new User with a connection
func NewConnectedUser(connection drivers.Connection) *User {
	user := NewUser("")
	user.conn = connection
	return user
}

// Listen starts reading from and writing to a user
func (u *User) Listen(chatServer *Server) {
	go u.ReadFrom(chatServer)
	go u.WriteTo()
}

// ReadFrom reads data from users and lets chat server interpret it
func (u *User) ReadFrom(chatServer *Server) {
	for {
		message := make([]byte, ReadConnectionLimitBytes)
		if _, err := u.conn.Read(message); err != nil {
			logs.ErrIfErrf(err, "Error reading from connection. @%s.", u.GetNickName())
			continue
		}

		_, err := u.HandleNewInput(chatServer, sanitizeInput(message))
		logs.ErrIfErrf(err, "Error handling input from user @%s.", u.GetNickName())
	}
}

func sanitizeInput(message []byte) string {
	message = bytes.Trim(message, "\x00")
	input := string(message)
	if strings.Contains(input, "\n") == true {
		input = strings.TrimSpace(input)
	}
	return input
}

// WriteTo to the user's connection and remembers the last message that was sent out
func (u *User) WriteTo() {
	for message := range u.outgoing {
		u.conn.Write([]byte(message + "\n"))
	}
}

// Disconnect a user from this server
func (u *User) Disconnect() error {
	nickName := u.GetNickName()
	logs.Infof("Disconnecting @%s", nickName)

	u.SetOutgoingf(plugins.UserOutPutTUserServerMessage, "Good Bye %f, come back again.", nickName)

	// Wait 1 second before actually disconnecting
	<-time.After(time.Second * 1)

	return u.conn.Close()
}
