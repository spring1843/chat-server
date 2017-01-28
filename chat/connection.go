package chat

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Connection is an interface for a network connection
type (
	MockedChatConnection struct {
		Outgoing  []byte
		Incoming  []byte
		Lock      *sync.Mutex
		EnableLog bool
	}
	Connection interface {
		Read(p []byte) (n int, err error)
		Write(p []byte) (n int, err error)
		Close() error
		RemoteAddr() net.Addr
	}
	MockedChatConnectionNetwork struct{}
)

func NewMockedChatConnection() *MockedChatConnection {
	fmt.Printf("Creating new mocked connection\n")
	return &MockedChatConnection{
		Lock:      new(sync.Mutex),
		Incoming:  make([]byte, 0),
		Outgoing:  make([]byte, 0),
		EnableLog: false,
	}
}

func (f *MockedChatConnectionNetwork) Network() string {
	return ``
}

func (f *MockedChatConnectionNetwork) String() string {
	return ``
}

func (m *MockedChatConnection) Write(p []byte) (int, error) {
	fmt.Printf("Writing data to connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Outgoing = p
	return len(m.Outgoing), nil
}

func (m *MockedChatConnection) ReadOutgoing() []byte {
	fmt.Printf("Reading outgoing data from connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	outgoing := m.Outgoing
	return outgoing
}

func (m *MockedChatConnection) Read(p []byte) (int, error) {
	fmt.Printf("Reading data from connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if len(m.Incoming) == 0 {
		return 0, io.EOF
	}
	i := 0
	for _, bit := range m.Incoming {
		p[i] = bit
		i++
	}
	return i, nil
}

func (m *MockedChatConnection) Close() error {
	fmt.Printf("Closing connection\n")
	return nil
}

func (m *MockedChatConnection) RemoteAddr() net.Addr {
	fmt.Printf("Reading remote address\n")
	return new(MockedChatConnectionNetwork)
}
