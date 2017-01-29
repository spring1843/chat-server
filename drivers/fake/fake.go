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
		Outgoing     []byte
		LockOutGoing *sync.Mutex

		Incoming     []byte
		LockIncoming *sync.Mutex

		EnableLog bool
		Lock      *sync.Mutex
	}
	// FakeNetwork is needed to implement the connection interface
	FakeNetwork struct{}
)

// NewFakeConnection returns a new fake connection
func NewFakeConnection() *FakeConnection {
	return &FakeConnection{
		Incoming:     make([]byte, 0),
		Outgoing:     make([]byte, 0),
		EnableLog:    false,
		LockOutGoing: new(sync.Mutex),
		LockIncoming: new(sync.Mutex),
		Lock:         new(sync.Mutex),
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
	conn.LockOutGoing.Lock()
	defer conn.LockOutGoing.Unlock()
	conn.Outgoing = p
	conn.log("End - Writing data to connection\n")
	return len(conn.Outgoing), nil
}

// ReadOutgoing reads what's going out to the user
func (conn *FakeConnection) ReadOutgoing() []byte {
	conn.log("Start - Reading outgoing data from connection\n")
	conn.LockOutGoing.Lock()
	defer conn.LockOutGoing.Unlock()
	outgoing := conn.Outgoing
	conn.log("End - Reading outgoing data from connection\n")
	return outgoing
}

// Read reads from the connection
func (conn *FakeConnection) Read(p []byte) (int, error) {
	conn.log("Start - Reading data from connection\n")
	conn.LockIncoming.Lock()
	defer conn.LockIncoming.Unlock()
	if len(conn.Incoming) == 0 {
		conn.log("End - EOF Reading data from connection\n")
		return 0, io.EOF
	}
	i := 0
	for _, bit := range conn.Incoming {
		p[i] = bit
		i++
	}
	conn.log("End - Reading data from connection\n")
	return i, nil
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
