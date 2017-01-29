package fake

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type (
	// FakeConnection is a chat server compatible connection used for testing
	FakeConnection struct {
		Outgoing  []byte
		Incoming  []byte
		Lock      *sync.Mutex
		EnableLog bool
	}
	// FakeNetwork is needed to implement the connection interface
	FakeNetwork struct{}
)

// NewFakeConnection returns a new fake connection
func NewFakeConnection() *FakeConnection {
	fmt.Printf("Creating new mocked connection\n")
	return &FakeConnection{
		Lock:      new(sync.Mutex),
		Incoming:  make([]byte, 0),
		Outgoing:  make([]byte, 0),
		EnableLog: false,
	}
}

// Write writes to connection
func (m *FakeConnection) Write(p []byte) (int, error) {
	fmt.Printf("Start - Writing data to connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Outgoing = p
	fmt.Printf("End - Writing data to connection\n")
	return len(m.Outgoing), nil
}

// ReadOutgoing reads what's going out to the user
func (m *FakeConnection) ReadOutgoing() []byte {
	fmt.Printf("Start - Reading outgoing data from connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	outgoing := m.Outgoing
	fmt.Printf("End - Reading outgoing data from connection\n")
	return outgoing
}

// Read reads from the connection
func (m *FakeConnection) Read(p []byte) (int, error) {
	fmt.Printf("Start - Reading data from connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if len(m.Incoming) == 0 {
		fmt.Printf("End - EOF Reading data from connection\n")
		return 0, io.EOF
	}
	i := 0
	for _, bit := range m.Incoming {
		p[i] = bit
		i++
	}
	fmt.Printf("End - Reading data from connection\n")
	return i, nil
}

// Close closes the connection
func (m *FakeConnection) Close() error {
	fmt.Printf("Closing connection\n")
	return nil
}

// RemoteAddr will return IP of client
func (m *FakeConnection) RemoteAddr() net.Addr {
	fmt.Printf("Reading remote address\n")
	return new(FakeNetwork)
}

// Network returns connection's origin network as string
func (f *FakeNetwork) Network() string {
	fmt.Printf("Reading network\n")
	return ``
}

// String returns network name as string
func (f *FakeNetwork) String() string {
	fmt.Printf("Returning string\n")
	return ``
}
