package chat

import (
	"io"
	"net"
	"sync"
)

// Connection is an interface for a network connection
type (
	MockedChatConnection struct {
		Outgoing []byte
		Incoming []byte
		Lock     *sync.Mutex
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
	return &MockedChatConnection{
		Lock: new(sync.Mutex),
	}
}

func (f *MockedChatConnectionNetwork) Network() string {
	return ``
}

func (f *MockedChatConnectionNetwork) String() string {
	return ``
}

func (m *MockedChatConnection) Write(p []byte) (n int, err error) {
	m.Lock.Lock()
	m.Outgoing = p
	m.Lock.Unlock()
	return len(m.Outgoing), nil
}

func (m *MockedChatConnection) ReadOutgoing() []byte {
	m.Lock.Lock()
	outgoing := m.Outgoing
	m.Lock.Unlock()
	return outgoing
}

func (m *MockedChatConnection) Read(p []byte) (n int, err error) {
	m.Lock.Lock()
	if len(m.Incoming) == 0 {
		m.Lock.Unlock()
		return 0, io.EOF
	}
	i := 0
	for _, bit := range m.Incoming {
		p[i] = bit
		i++
	}
	m.Lock.Unlock()
	return i, nil
}

func (m *MockedChatConnection) Close() error {
	return nil
}

func (m *MockedChatConnection) RemoteAddr() net.Addr {
	return new(MockedChatConnectionNetwork)
}
