package websocket

import (
	"net"

	"time"

	"github.com/spring1843/chat-server/libs/websocket"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// ChatConnection is an middleman between the WebSocket connection and Chat Server
type ChatConnection struct {
	Connection *websocket.Conn
	incoming   chan []byte
	outgoing   chan []byte
}

// NewChatConnection returns a new ChatConnection
func NewChatConnection() *ChatConnection {
	return &ChatConnection{
		incoming: make(chan []byte),
		outgoing: make(chan []byte),
	}
}

// Read waits for user to enter a text, or reads the last entered incoming message
func (c *ChatConnection) Read(p []byte) (int, error) {
	i := 0
	message := <-c.incoming
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

// Write to a ChatConnection
func (c *ChatConnection) Write(p []byte) (int, error) {
	// At max 1 go routine must be writing to this connection so we use a channel here
	c.outgoing <- p
	return len(p), nil

}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ChatConnection) readPump() {
	defer func() {
		c.Connection.Close()
		logs.Infof("No longer reading Websocket pump for %s", c.Connection.RemoteAddr())
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
		c.incoming <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ChatConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()
	for {
		select {
		case message, ok := <-c.outgoing:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.outgoing)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.outgoing)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
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
