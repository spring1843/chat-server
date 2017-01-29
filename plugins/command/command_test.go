package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/plugins/command"
)

var (
	user1 = chat.NewUser("u1")
	user2 = chat.NewUser("u2")
	user3 = chat.NewUser("u3")
)

func Test_CanValidate(t *testing.T) {
	var (
		invalidCommand1 = ``
		invalidCommand2 = `badcommand`
		invalidCommand3 = `/badcommand`
		validCommand1   = `/help`
		validCommand2   = `/join`
	)

	if _, err := command.GetCommand(invalidCommand1); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand1)
	}

	if _, err := command.GetCommand(invalidCommand2); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand2)
	}

	if _, err := command.GetCommand(invalidCommand3); err == nil {
		t.Errorf("Invalid command was detected valid, got %s", invalidCommand3)
	}

	if _, err := command.GetCommand(validCommand1); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand1)
	}

	if _, err := command.GetCommand(validCommand2); err != nil {
		t.Errorf("Valid command was detected invalid, got %s", validCommand2)
	}
}

func Test_HelpCommand(t *testing.T) {
	fakeConnection := chat.NewFakeConnection()
	fakeConnection.Incoming = []byte("/help\n")

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	server.AddUser(user)
	msg := user.GetOutgoing()

	if strings.Contains(msg, "Shows the list of all available commands") != true {
		t.Errorf("Help command did not output description of help command")
	}
}

func Test_ListCommand(t *testing.T) {
	fakeConnection1 := chat.NewFakeConnection()
	fakeConnection2 := chat.NewFakeConnection()
	fakeConnection3 := chat.NewFakeConnection()

	server := chat.NewServer()
	server.AddChannel(`foo`)

	input := `/join #foo`
	joinCommand, err := command.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user1.SetNickName(`u1`)

	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user1.SetNickName(`u2`)

	user3 := chat.NewConnectedUser(server, fakeConnection3)
	user1.SetNickName(`u3`)

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)

	user1.ExecuteCommand(server, input, joinCommand)
	user2.ExecuteCommand(server, input, joinCommand)
	user3.ExecuteCommand(server, input, joinCommand)

	input = "/list \n"
	listCommand, err := command.GetCommand(input)
	if err != nil {
		t.Errorf("Could not get an instance of list command")
	}
	user1.ExecuteCommand(server, input, listCommand)
	msg := user1.GetOutgoing()

	if strings.Contains(msg, "@u1") != true && strings.Contains(msg, "@u2") != true && strings.Contains(msg, "@u3") != true {
		t.Errorf("List command did not show all users in the room")
	}
}

func Test_IgnoreCommand(t *testing.T) {
	fakeConnection1 := chat.NewFakeConnection()
	fakeConnection2 := chat.NewFakeConnection()

	server := chat.NewServer()

	user1 := chat.NewConnectedUser(server, fakeConnection1)
	user2 := chat.NewConnectedUser(server, fakeConnection2)
	user2.SetNickName(`u2`)

	server.AddUser(user1)
	server.AddUser(user2)

	input := `/ignore @u2`
	ignoreCommand, _ := command.GetCommand(input)
	user1.ExecuteCommand(server, input, ignoreCommand)
	if user1.HasIgnored(user2.GetNickName()) != true {
		t.Errorf("User was not ignored after ignore command executed")
	}
}

func Test_JoinCommand(t *testing.T) {
	fakeConnection := chat.NewFakeConnection()
	fakeConnection.Incoming = []byte("/join #r\n")

	server := chat.NewServer()
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
	server := chat.NewServer()

	server.AddUser(user1)
	server.AddUser(user2)

	channel := chat.NewChannel()
	channel.Name = `r`
	user1.SetChannel(channel.Name)
	user2.SetChannel(channel.Name)
	input := `/msg @u2 foo`
	messageCommand, err := command.GetCommand(input)
	if err != nil {
		t.Fatalf("Failed getting message command. Error %s", err)
	}

	if err := user1.ExecuteCommand(server, input, messageCommand); err != nil {
		t.Fatalf("Failed executing message. Error %s", err)
	}
	msg := user2.GetOutgoing()

	if strings.Contains(msg, "- *Private from @u1: foo") != true {
		t.Errorf("Private message was not received. Last message %s", msg)
	}
}

func Test_QuitCommand(t *testing.T) {
	fakeConnection := chat.NewFakeConnection()

	server := chat.NewServer()
	user := chat.NewConnectedUser(server, fakeConnection)
	user.SetNickName(`foo`)
	server.AddUser(user)

	if server.IsUserConnected(`foo`) != true {
		t.Errorf("User was  disconnected without runnign the quit command")
	}

	input := `/quit`
	quitCommand, _ := command.GetCommand(input)
	user.ExecuteCommand(server, input, quitCommand)

	if server.IsUserConnected(`foo`) {
		t.Errorf("User was not disconnected after running quit command")
	}
}
