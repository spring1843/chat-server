package fake

import (
	"io"
	"net"
	"sync"

	"github.com/spring1843/chat-server/plugins/logs"
)

type (
	// FakeConnection is a chat server compatible connection used for testing
	FakeConnection struct {
		outgoing     []byte
		lockOutGoing *sync.Mutex

		incoming     []byte
		lockIncoming *sync.Mutex

		EnableLog bool
	}
	// FakeNetwork is needed to implement the connection interface
	FakeNetwork struct{}
)

// NewFakeConnection returns a new fake connection
func NewFakeConnection() *FakeConnection {
	return &FakeConnection{
		incoming:     make([]byte, 0),
		outgoing:     make([]byte, 0),
		EnableLog:    false,
		lockOutGoing: new(sync.Mutex),
		lockIncoming: new(sync.Mutex),
	}
}

const logPrefix = "Fake Connection"

func (m *FakeConnection) log(message string) {
	if m.EnableLog {
		logs.Infof("%s - %s", logPrefix, message)
	}
}

// Write writes to connection
func (conn *FakeConnection) Write(p []byte) (int, error) {
	conn.log("Start - Writing data to connection\n")
	conn.lockOutGoing.Lock()
	defer conn.lockOutGoing.Unlock()
	conn.outgoing = p
	conn.log("End - Writing data to connection\n")
	return len(conn.outgoing), nil
}

// ReadOutgoing reads what's going out to the user
func (conn *FakeConnection) ReadOutgoing() []byte {
	conn.log("Start - Reading outgoing data from connection\n")
	conn.lockOutGoing.Lock()
	defer conn.lockOutGoing.Unlock()
	outgoing := conn.outgoing
	conn.log("End - Reading outgoing data from connection\n")
	return outgoing
}

// Read reads from the connection
func (conn *FakeConnection) Read(p []byte) (int, error) {
	conn.log("Start - Reading data from connection\n")
	conn.lockIncoming.Lock()
	defer conn.lockIncoming.Unlock()
	if len(conn.incoming) == 0 {
		conn.log("End - EOF Reading data from connection\n")
		return 0, io.EOF
	}
	i := 0
	for _, bit := range conn.incoming {
		p[i] = bit
		i++
	}
	conn.log("End - Reading data from connection\n")
	return i, nil
}

// GetOutgoing gets the outgoing message for a user
func (conn *FakeConnection) GetOutgoing() string {
	conn.lockOutGoing.Lock()
	defer conn.lockOutGoing.Unlock()
	return string(conn.outgoing)
}

// SetOutgoing sets an outgoing message to the user
func (conn *FakeConnection) SetOutgoing(message string) {
	conn.lockOutGoing.Lock()
	defer conn.lockOutGoing.Unlock()
	conn.outgoing = []byte(message)
}

// GetIncoming gets the incoming message from the user
func (conn *FakeConnection) GetIncoming() string {
	conn.lockIncoming.Lock()
	defer conn.lockIncoming.Unlock()
	return string(conn.incoming)
}

// SetIncoming sets an incoming message from the user
func (conn *FakeConnection) SetIncoming(message string) {
	conn.lockIncoming.Lock()
	defer conn.lockIncoming.Unlock()
	conn.incoming = []byte(message)
}

// Close closes the connection
func (conn *FakeConnection) Close() error {
	conn.log("Closing connection\n")
	return nil
}

// RemoteAddr will return IP of client
func (conn *FakeConnection) RemoteAddr() net.Addr {
	conn.log("Reading remote address\n")
	return new(FakeNetwork)
}

// Network returns connection's origin network as string
func (f *FakeNetwork) Network() string {
	return ``
}

// String returns network name as string
func (f *FakeNetwork) String() string {
	return ``
}
