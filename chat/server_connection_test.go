package chat_test

import (
	"os"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_WelcomeNewUsers(t *testing.T) {
	server = chat.NewServer()
	server.Listen()

	conn1 := chat.NewFakeConnection()
	conn1.Lock.Lock()
	defer conn1.Lock.Unlock()
	conn1.Incoming = []byte("foo\n")

	server.WelcomeNewUser(conn1)
	if !server.IsUserConnected("foo") {
		t.Error("User foo not added to the server")
	}

	connection2 := chat.NewFakeConnection()
	connection2.Lock.Lock()
	defer connection2.Lock.Unlock()
	connection2.Incoming = []byte("bar\n")

	server.WelcomeNewUser(connection2)
	if !server.IsUserConnected("bar") {
		t.Error("User bar not added to the server")
	}
}
