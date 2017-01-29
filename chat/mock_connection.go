package chat

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type (
	// FakeConnection is an interface for a network connection
	FakeConnection struct {
		Outgoing  []byte
		Incoming  []byte
		Lock      *sync.Mutex
		EnableLog bool
	}
	FakeNetwork struct{}
)

func NewMockedChatConnection() *FakeConnection {
	fmt.Printf("Creating new mocked connection\n")
	return &FakeConnection{
		Lock:      new(sync.Mutex),
		Incoming:  make([]byte, 0),
		Outgoing:  make([]byte, 0),
		EnableLog: false,
	}
}

func (m *FakeConnection) Write(p []byte) (int, error) {
	fmt.Printf("Start - Writing data to connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	m.Outgoing = p
	fmt.Printf("End - Writing data to connection\n")
	return len(m.Outgoing), nil
}

func (m *FakeConnection) ReadOutgoing() []byte {
	fmt.Printf("Start - Reading outgoing data from connection\n")
	m.Lock.Lock()
	defer m.Lock.Unlock()
	outgoing := m.Outgoing
	fmt.Printf("End - Reading outgoing data from connection\n")
	return outgoing
}

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

func (f *FakeNetwork) Network() string {
	fmt.Printf("Reading network\n")
	return ``
}

func (f *FakeNetwork) String() string {
	fmt.Printf("Returning string\n")
	return ``
}

func (m *FakeConnection) Close() error {
	fmt.Printf("Closing connection\n")
	return nil
}

func (m *FakeConnection) RemoteAddr() net.Addr {
	fmt.Printf("Reading remote address\n")
	return new(FakeNetwork)
}
