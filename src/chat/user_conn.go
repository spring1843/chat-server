package chat

import (
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/spring1843/chat-server/src/drivers"
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

// GetOutgoing gets the outgoing message for a user
func (u *User) GetOutgoing() string {
	return <-u.outgoing
}

// SetOutgoing sets an outgoing message to the user
func (u *User) SetOutgoing(message string) {
	u.outgoing <- message
}

// GetIncoming gets the incoming message from the user
func (u *User) GetIncoming() string {
	return <-u.incoming
}

// SetIncoming sets an incoming message from the user
func (u *User) SetIncoming(message string) {
	u.incoming <- message
}

// ReadFrom reads data from users and lets chat server interpret it
func (u *User) ReadFrom(chatServer *Server) {
	for {
		message := make([]byte, ReadConnectionLimitBytes)
		if _, err := u.conn.Read(message); err != nil {
			if err == io.EOF {
				continue
			}
			logs.ErrIfErrf(err, "Error reading from @%s.", u.GetNickName())
		}

		message = bytes.Trim(message, "\x00")

		input := string(message)
		//Remove new line
		if strings.Contains(input, "\n") == true {
			input = strings.TrimSpace(input)
		}

		handled, err := u.HandleNewInput(chatServer, input)
		if err != nil {
			logs.ErrIfErrf(err, "Error reading input from user @%s.", u.nickName)
		}
		if handled {
			//If handled then continue reading
			continue
		}

		if input != "\n" && input != `` {
			u.SetIncoming(input)
		}
	}
}

// WriteTo to the user's connection and remembers the last message that was sent out
func (u *User) WriteTo() {
	for message := range u.outgoing {
		u.conn.Write([]byte(message + "\n"))
	}
}

// Disconnect a user from this server
func (u *User) Disconnect() error {
	logs.Infof("disconnecting=@%s", u.nickName)
	u.SetOutgoing("Good Bye, come back again.")

	// Wait 1 second before actually disconnecting
	<-time.After(time.Second * 1)
	return u.conn.Close()
}
