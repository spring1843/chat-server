package websocket

import (
	"net"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

// ChatConnection is an middleman between the WebSocket connection and Chat Server
type ChatConnection struct {
	Connection *websocket.Conn
	Incoming   chan []byte
}

// NewChatConnection returns a new ChatConnection
func NewChatConnection() *ChatConnection {
	return &ChatConnection{
		Incoming: make(chan []byte),
	}
}

// Write to a ChatConnection
func (c *ChatConnection) Write(p []byte) (int, error) {
	w, err := c.Connection.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, errs.Wrap(err, "Error getting nextwriter from WebSocket connection.")
	}
	defer w.Close()
	return w.Write(p)
}

// Close a ChatConnection
func (c *ChatConnection) Close() error {
	if err := c.Connection.Close(); err != nil {
		return errs.Wrap(err, "Error closing WebSocket connection")
	}
	return nil
}

// RemoteAddr returns the remote address of the connected user
func (c *ChatConnection) RemoteAddr() net.Addr {
	return c.Connection.RemoteAddr()
}

// Read from a ChatConnection
// P is a buffered, write only from the start and maintain the size
func (c *ChatConnection) Read(p []byte) (int, error) {
	i := 0
	message := <-c.Incoming
	message = append(message, byte('\n'))

	if len(p) < len(message) {
		p = make([]byte, len(message), len(message))
	}

	for i, bit := range message {
		p[i] = bit
	}
	return i, nil
}

func listen(c *ChatConnection) {
	for {
		msgType, message, err := c.Connection.ReadMessage()
		if err != nil {
			logs.ErrIfErrf(err, "Error reading from WebSocket connection")
			break
		}
		if msgType == 1 {
			c.Incoming <- message
		}
	}
	logs.Infof("No longer listening to %s", c.RemoteAddr())
}
