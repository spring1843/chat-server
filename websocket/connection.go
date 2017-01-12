package websocket

import (
	"net"

	"github.com/gorilla/websocket"
)

// ChatConnection is an middleman between the WebSocket connection and Chat Server
type ChatConnection struct {
	Connection *websocket.Conn
	Incoming   chan []byte
}

// NewChatConnection returns a new ChatConnection
func NewChatConnection() *ChatConnection {
	newChatConnection := &ChatConnection{
		Incoming: make(chan []byte),
	}
	return newChatConnection
}

// Write to a ChatConnection
func (c *ChatConnection) Write(p []byte) (int, error) {
	err := handleOutgoing(1, c, p)
	if err != nil {
		return -1, err
	}
	return len(p) - 1, nil
}

// Close a ChatConnection
func (c *ChatConnection) Close() error {
	err := c.Connection.Close()
	if err != nil {
		return err
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
		p = make([]byte, len(message))
	}

	for _, bit := range message {
		p[i] = bit
		i++
	}
	return i, nil
}

func handleIncoming(c *ChatConnection) error {
	msgType, message, err := c.Connection.ReadMessage()
	if err != nil {
		return err
	}
	if msgType == 1 {
		c.Incoming <- message
	}
	return nil
}

func handleOutgoing(msgType int, c *ChatConnection, message []byte) error {
	err := c.Connection.WriteMessage(msgType, message)
	if err != nil {
		return err
	}
	return nil
}

func listen(c *ChatConnection) {
	for {
		err := handleIncoming(c)
		if err != nil {
			c.Close()
			break
		}
	}
}
