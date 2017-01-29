package fake

import (
	"io"
	"net"
	"sync"

	"github.com/spring1843/chat-server/plugins/errs"
	"github.com/spring1843/chat-server/plugins/logs"
)

const maxDataLength = 255

type (
	// FakeConnection is a chat server compatible connection used for testing
	FakeConnection struct {
		data      []byte
		lock      *sync.Mutex
		EnableLog bool
	}
	// FakeNetwork is needed to implement the connection interface
	FakeNetwork struct{}
)

// NewFakeConnection returns a new fake connection
func NewFakeConnection() *FakeConnection {
	return &FakeConnection{
		data:      make([]byte, 0),
		EnableLog: false,
		lock:      new(sync.Mutex),
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
		if i >= len(p) {
			return 0, errs.Newf("Input too small. Length of data is %d", len(data))
		}
		p[i] = bit
	}

	conn.log("End - Reading data from connection\n")
	return len(data), nil
}

// ReadString convenient method for reading
func (conn *FakeConnection) ReadString(length int) (string, error) {
	data := make([]byte, length, length)
	n, err := conn.Read(data)
	if err != nil {
		return "", errs.Wrap(err, "Error reading from fake connection.")
	}
	if n == 0 {
		return "", errs.New("Error reading from fake connection. Length of data is zero")
	}
	return string(data), nil
}

// WriteString convenient method for writing
func (conn *FakeConnection) WriteString(message string) (int, error) {
	return conn.Write([]byte(message))
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
