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
	user := chat.NewConnectedUser(server, fakeConnection)
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

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user1.NickName = `u1`

	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user2.NickName = `u2`

	user3 := chat.NewConnectedUser(server, fakeConnection3)
	user3.NickName = `u3`

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	user1.ExecuteCommand(server, input, joinCommand)
	user2.ExecuteCommand(server, input, joinCommand)
	user3.ExecuteCommand(server, input, joinCommand)

	input = "/list \n"
	listCommand, err := chat.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}
	user1.ExecuteCommand(server, input, listCommand)
	msg := user1.LastOutGoingMessage

	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}

func Test_IgnoreCommand(t *testing.T) {
	fakeConnection1 := NewMockedChatConnection()
	fakeConnection2 := NewMockedChatConnection()

	server := chat.NewService()

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user2.NickName = `u2`

	server.AddUser(user1)
	server.AddUser(user2)

	input := `/ignore @u2`
	ignoreCommand, _ := chat.GetCommand(input)
	user1.ExecuteCommand(server, input, ignoreCommand)
	if user1.HasIgnored(user2.NickName) != true {
		t.Errorf("User was not ignored after ignore command executed")
	}
}

func Test_JoinCommand(t *testing.T) {
	fakeConnection := NewMockedChatConnection()
	fakeConnection.incoming = []byte("/join #r\n")

	server := chat.NewService()
	user1 := chat.NewConnectedUser(server, fakeConnection)
	server.AddUser(user1)

	msg := user1.GetOutgoing()
	if strings.Contains(msg, "other users this channel") != true {
		t.Errorf("User did not receive welcome message after joining channel received %s instead", msg)
	}

	if user1.GetChannel() != "" && user1.GetChannel() != `r` {
		t.Errorf("User did not join the #r channel when he should have")
	}
}

func Test_MessageCommand(t *testing.T) {
	server := chat.NewService()

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.Name = `r`
	user1.SetChannel(channel.Name)
	user2.SetChannel(channel.Name)
	input := `/msg @u2 foo`
	messageCommand, err := chat.GetCommand(input)
	if err != nil {
		t.Fatalf("Failed getting message command. Error %s", err)
	}

	if err := user1.ExecuteCommand(server, input, messageCommand); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}
	msg := user2.GetOutgoing()

	if strings.Contains(msg, "- *Private from @u1: foo") != true {
		t.Errorf("Private message was not received. Last message %s", user2.LastOutGoingMessage)
	}
}

func Test_QuitCommand(t *testing.T) {
	fakeConnection := NewMockedChatConnection()

	server := chat.NewService()
	user := chat.NewConnectedUser(server, fakeConnection)
	user.NickName = `foo`
	server.AddUser(user)

	if server.IsUserConnected(`foo`) != true {
		t.Errorf("User was  disconnected without runnign the quit command")
	}

	input := `/quit`
	quitCommand, _ := chat.GetCommand(input)
	user.ExecuteCommand(server, input, quitCommand)

	if server.IsUserConnected(`foo`) {
		t.Errorf("User was not disconnected after running quit command")
	}
}
