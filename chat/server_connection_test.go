package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func Test_WelcomeNewUsers(t *testing.T) {
	var (
		server     = chat.NewServer()
		connection = fake.NewFakeConnection()
	)

	server.Listen()

	connection.Lock.Lock()
	connection.Incoming = []byte("foo\n")
	connection.Lock.Unlock()

	server.ConnectUser(connection)

	if len(server.Users) != 1 {
		t.Errorf("User was not added to the server")
	}
}
