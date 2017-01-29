package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/drivers/fake"
)

func TestInterviewUser(t *testing.T) {

	t.Skipf("Racy!")

	var (
		server     = chat.NewServer()
		connection = fake.NewFakeConnection()
	)

	server.Listen()

	connection.LockIncoming.Lock()
	connection.Incoming = []byte("newuser\n")
	connection.LockIncoming.Unlock()

	server.InterviewUser(connection)
	if server.ConnectedUsersCount() != 1 {
		t.Errorf("User was not added to the server")
	}
}
