// +build !race

package chat_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
)

func Test_CanValidate(t *testing.T) {
	var (
		invalidCommand1 = ``
		invalidCommand2 = `badcommand`
		invalidCommand3 = `/badcommand`
		validCommand1   = `/help`
		validCommand2   = `/join`
	)

	if _, err := chat.GetCommand(invalidCommand1); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand1)
	}

	if _, err := chat.GetCommand(invalidCommand2); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand2)
	}

	if _, err := chat.GetCommand(invalidCommand3); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand3)
	}

	if _, err := chat.GetCommand(validCommand1); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand1)
	}

	if _, err := chat.GetCommand(validCommand2); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand2)
	}
}

func Test_HelpCommand(t *testing.T) {
	fakeConnection := NewMockedChatConnection()
	fakeConnection.incoming = []byte("/help\n")

	server := chat.NewService()
	user := chat.NewConnectedUser(fakeConnection)
	server.AddUser(user)
	msg := user.GetOutgoing()

	if strings.Contains(msg, "Shows the list of all available commands") != true {
		t.Errorf("Help command did not output description of help command")
	}
}

func Test_ListCommand(t *testing.T) {
	fakeConnection1 := NewMockedChatConnection()
	fakeConnection2 := NewMockedChatConnection()
	fakeConnection3 := NewMockedChatConnection()

	server := chat.NewService()
	server.AddChannel(`foo`)

	input := `/join #foo`
	joinCommand, err := chat.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}

	user1 := chat.NewConnectedUser(fakeConnection1)
	user1.NickName = `u1`

	user2 := chat.NewConnectedUser(fakeConnection2)
	user2.NickName = `u2`

	user3 := chat.NewConnectedUser(fakeConnection3)
	user3.NickName = `u3`

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	user1.ExecuteCommand(input, joinCommand)
	user2.ExecuteCommand(input, joinCommand)
	user3.ExecuteCommand(input, joinCommand)

	input = "/list \n"
	listCommand, err := chat.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}
	user1.ExecuteCommand(input, listCommand)
	msg := user1.LastOutGoingMessage

	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}

func Test_IgnoreCommand(t *testing.T) {
	fakeConnection1 := NewMockedChatConnection()
	fakeConnection2 := NewMockedChatConnection()

	server := chat.NewService()

	user1 := chat.NewConnectedUser(fakeConnection1)
	user2 := chat.NewConnectedUser(fakeConnection2)
	user2.NickName = `u2`

	server.AddUser(user1)
	server.AddUser(user2)

	input := `/ignore @u2`
	ignoreCommand, _ := chat.GetCommand(input)
	user1.ExecuteCommand(input, ignoreCommand)
	if user1.HasIgnored(user2.NickName) != true {
		t.Errorf("User was not ignored after ignore command executed")
	}
}

func Test_JoinCommand(t *testing.T) {
	fakeConnection1 := NewMockedChatConnection()

	server := chat.NewService()
	user1 := chat.NewConnectedUser(fakeConnection1)
	server.AddUser(user1)

	input := `/join #r`
	joinCommand, err := chat.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}

	user1.ExecuteCommand(input, joinCommand)
	if user1.Channel != nil && user1.Channel.Name != `r` {
		t.Errorf("User did not join the #r channel when he should have")
	}
}

func Test_MessageCommand(t *testing.T) {
	fakeConnection1 := NewMockedChatConnection()
	fakeConnection2 := NewMockedChatConnection()

	server := chat.NewService()

	user1 := chat.NewConnectedUser(fakeConnection1)
	user1.NickName = `u1`

	user2 := chat.NewConnectedUser(fakeConnection2)
	user2.NickName = `u2`

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.Name = `r`
	user1.Channel, user2.Channel = channel, channel

	input := `/msg @u2 foo`
	messageCommand, _ := chat.GetCommand(input)
	user1.ExecuteCommand(input, messageCommand)

	msg := user2.LastOutGoingMessage

	if strings.Contains(msg, "- *Private from @u1: foo") != true {
		t.Errorf("Private message was not received")
	}
}

func Test_QuitCommand(t *testing.T) {
	fakeConnection := NewMockedChatConnection()

	server := chat.NewService()
	user := chat.NewConnectedUser(fakeConnection)
	user.NickName = `foo`
	server.AddUser(user)

	if server.IsUserConnected(`foo`) != true {
		t.Errorf("User was  disconnected without runnign the quit command")
	}

	input := `/quit`
	quitCommand, _ := chat.GetCommand(input)
	user.ExecuteCommand(input, quitCommand)

	if server.IsUserConnected(`foo`) != false {
		t.Errorf("User was not disconnected after running quit command")
	}
}
