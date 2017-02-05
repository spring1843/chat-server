package chat_test

import (
	"fmt"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

func TestUserInterview(t *testing.T) {
	tryouts := 10
	server := chat.NewServer()
	server.Listen()

	connections := make([]*fake.MockedConnection, tryouts, tryouts)
	i := 0
	for i < tryouts {
		connections[i] = fake.NewFakeConnection()

		nickName := fmt.Sprintf("user%d", i)
		n, err := connections[i].WriteString(nickName)
		if err != nil {
			t.Fatalf("Failed writing to connection for user %q. Error %s", nickName, err)
		}
		if n != len(nickName) {
			t.Fatalf("Wrong length after write. For user %s Expected %d, got %d.", nickName, len(nickName), n)
		}
		server.InterviewUser(connections[i])
		if server.ConnectedUsersCount() != i+1 {
			t.Errorf("User %s was not added to the server", nickName)
		}
		i++
	}
}
