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


		data     []byte
		lock *sync.Mutex
		EnableLog bool
	}
	// FakeNetwork is needed to implement the connection interface
	FakeNetwork struct{}
)

// NewFakeConnection returns a new fake connection
func NewFakeConnection() *FakeConnection {
	return &FakeConnection{
		data:     make([]byte, 0),
		incoming:     make([]byte, 0),
		outgoing:     make([]byte, 0),
		EnableLog:    false,
		lock: new(sync.Mutex),
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
	conn.log("Lock\n")
	conn.lock.Lock()
	conn.log("Start - Writing data to connection\n")

	conn.data = p
	n := len(conn.data)

	conn.log("Unlock\n")
	conn.lock.Unlock()
	conn.log("End - Writing data to connection\n")

	return n, nil
}

// Read reads from the connection
func (conn *FakeConnection) Read(p []byte) (int, error) {
	conn.log("Lock\n")
	conn.lock.Lock()
	data := conn.data
	conn.log("Unlock\n")
	conn.lock.Unlock()

	conn.log("Start - Reading data from connection\n")
	if len(data) == 0 {
		conn.log("End - EOF Reading data from connection\n")
		return 0, io.EOF
	}

	for i, bit := range data {
		p[i] = bit
	}

	conn.log("End - Reading data from connection\n")
	return len(data), nil
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
