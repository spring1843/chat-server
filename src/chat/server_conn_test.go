package chat_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestInterviewUser(t *testing.T) {
	var (
		server     = chat.NewServer()
		connection = fake.NewFakeConnection()
	)

	server.Listen()

	input := "newuser\n"
	n, err := connection.WriteString(input)
	if err != nil {
		t.Fatalf("Failed writing to connection. Error %s", err)
	}
	if n != len(input) {
		t.Fatalf("Wrong length after write. Expected %d, got %d.", len(input), n)
	}

	server.InterviewUser(connection)
	if server.ConnectedUsersCount() != 1 {
		t.Errorf("User was not added to the server")
	}
}
