package websocket

import (
	"net"

	"time"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const (
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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

// Write to a ChatConnection
func (c *ChatConnection) Write(p []byte) (int, error) {
	w, err := c.Connection.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, errs.Wrap(err, "Error getting nextwriter from WebSocket connection.")
	}
	defer w.Close()
	return w.Write(p)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ChatConnection) readPump() {
	defer func() {
		c.Connection.Close()
	}()
	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	c.Connection.SetPongHandler(func(string) error { c.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logs.ErrIfErrf(err, "Error reading from WebSocket connection")
			}
			break
		}
		c.Incoming <- message
	}
	logs.Infof("No longer listening to %s", c.RemoteAddr())
}

// RemoteAddr returns the remote address of the connected user
func (c *ChatConnection) RemoteAddr() net.Addr {
	return c.Connection.RemoteAddr()
}

// Close a ChatConnection
func (c *ChatConnection) Close() error {
	if err := c.Connection.Close(); err != nil {
		return errs.Wrap(err, "Error closing WebSocket connection")
	}
	return nil
}
