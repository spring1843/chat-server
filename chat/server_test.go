package chat_test

import (
	"os"
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

var (
	server = chat.NewServer()
)

func Test_CanLogToFile(t *testing.T) {
	fakeWriter := chat.NewMockedChatConnection()
	server.SetLogFile(fakeWriter)
	server.LogPrintf("test \t foo\n")

	logMessage := string(fakeWriter.Outgoing)

	if strings.Contains(logMessage, `foo`) == false {
		t.Errorf("Did not send log to file")
	}
}

func Test_CanAddUser(t *testing.T) {
	server.AddUser(user1)
	if !server.IsUserConnected(`u1`) {
		t.Errorf("User is not connected when should have been connected")
	}
	if server.IsUserConnected(`bar`) {
		t.Errorf("User is connected when should not have been connected")
	}
}

func Test_CanRemoveUser(t *testing.T) {
	server.AddUser(user1)
	server.AddUser(user2)

	server.RemoveUser(user1.GetNickName())

	if server.IsUserConnected(`u1`) {
		t.Errorf("User is was not removed when should have been")
	}

	if server.ConnectedUsersCount() != 1 {
		t.Errorf("After adding two users and removing one user total users does not equal 1")
	}
}

func Test_AddChannel(t *testing.T) {
	server.AddChannel(`foo`)

	if server.GetChannelCount() != 1 {
		t.Errorf("Couldn't add a channel")
	}
}

func Test_GetSameChannel(t *testing.T) {
	server.AddChannel(`foo`)
	sameChannel, err := server.GetChannel(`foo`)

	if err != nil || "foo" != sameChannel.GetName() {
		t.Errorf("Couldn't add and get channel")
	}
}

func Test_WelcomeNewUsers(t *testing.T) {
	server = chat.NewServer()
	logFile, _ := os.OpenFile(`/dev/null`, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	server.SetLogFile(logFile)

	server.Listen()

	connection1 := chat.NewMockedChatConnection()
	connection1.Lock.Lock()
	defer connection1.Lock.Unlock()
	connection1.Incoming = []byte("foo\n")

	server.WelcomeNewUser(connection1)
	if !server.IsUserConnected("foo") {
		t.Error("User foo not added to the server")
	}

	connection2 := chat.NewMockedChatConnection()
	connection2.Lock.Lock()
	defer connection2.Lock.Unlock()
	connection2.Incoming = []byte("bar\n")

	server.WelcomeNewUser(connection2)
	if !server.IsUserConnected("bar") {
		t.Error("User bar not added to the server")
	}
}
