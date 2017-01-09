package chat_test

import (
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_CanLogToFile(t *testing.T) {
	var server = chat.NewServer()

	fakeWriter := NewMockedChatConnection()
	server.SetLogFile(fakeWriter)
	server.LogPrintf("test \t foo\n")

	logMessage := string(fakeWriter.outgoing)

	if strings.Contains(logMessage, `foo`) == false {
		t.Errorf("Did not send log to file")
	}
}

func Test_CanAddUser(t *testing.T) {
	var (
		server = chat.NewServer()
		user   = new(chat.User)
	)

	user.NickName = `foo`
	server.AddUser(user)
	if server.IsUserConnected(`foo`) != true {
		t.Errorf("User is not connected when should have been connected")
	}
	if server.IsUserConnected(`bar`) != false {
		t.Errorf("User is connected when should not have been connected")
	}
}

func Test_CanRemoveUser(t *testing.T) {
	var (
		server = chat.NewServer()
		user1  = new(chat.User)
		user2  = new(chat.User)
	)

	user1.NickName = `u1`
	server.AddUser(user1)

	user2.NickName = `u2`
	server.AddUser(user2)

	err := server.RemoveUser(user1)

	if err != nil || server.IsUserConnected(`u1`) != false {
		t.Errorf("User is was not removed when should have been")
	}

	if len(server.Users) != 1 {
		t.Errorf("After adding two users and removing one user total users does not equal 1")
	}
}

func Test_AddChannel(t *testing.T) {
	var (
		server = chat.NewServer()
	)

	channel := server.AddChannel(`foo`)

	if server.Channels[0].Name != channel.Name {
		t.Errorf("Couldn't add a channel")
	}
}

func Test_GetSameChannel(t *testing.T) {
	var (
		server = chat.NewServer()
	)

	channel := server.AddChannel(`foo`)
	sameChannel, err := server.GetChannel(`foo`)

	if err != nil || channel != sameChannel {
		t.Errorf("Couldn't add and get channel")
	}
}

func Test_WelcomeNewUser(t *testing.T) {
	var (
		server     = chat.NewServer()
		connection = NewMockedChatConnection()
	)

	logFile, _ := os.OpenFile(`/dev/null`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	server.SetLogFile(logFile)

	server.Listen()
	connection.incomingMutex.Lock()
	connection.incoming = []byte("foo\n")
	connection.incomingMutex.Unlock()

	server.WelcomeNewUser(connection)

	if len(server.Users) != 1 {
		t.Errorf("User was not added to the server")
	}
}

type MockedChatConnectionNetwork struct{}

func NewMockedChatConnection() *MockedChatConnection {
	mockedChatConnection := new(MockedChatConnection)
	mockedChatConnection.incomingMutex = &sync.Mutex{}
	return mockedChatConnection
}

func (f *MockedChatConnectionNetwork) Network() string {
	return ``
}
func (f *MockedChatConnectionNetwork) String() string {
	return ``
}

type MockedChatConnection struct {
	outgoing      []byte
	incoming      []byte
	incomingMutex *sync.Mutex
}

func (m *MockedChatConnection) Write(p []byte) (n int, err error) {
	m.outgoing = p
	return len(m.outgoing), nil
}

func (m *MockedChatConnection) Read(p []byte) (n int, err error) {
	m.incomingMutex.Lock()
	if len(m.incoming) == 0 {
		m.incomingMutex.Unlock()
		return 0, io.EOF
	}
	i := 0
	for _, bit := range m.incoming {
		p[i] = bit
		i++
	}
	m.incomingMutex.Unlock()
	return i, nil
}

func (m *MockedChatConnection) Close() error {
	return nil
}

func (m *MockedChatConnection) RemoteAddr() net.Addr {
	return new(MockedChatConnectionNetwork)
}
