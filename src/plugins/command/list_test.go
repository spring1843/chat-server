package command_test

import (
	"strings"
	"testing"

	"github.com/spring1843/chat-server/src/chat"
)

func TestListCommand(t *testing.T) {
	server := chat.NewServer()
	user1, user2, user3, user4, user5 := chat.NewUser("u1"), chat.NewUser("u2"), chat.NewUser("u3"), chat.NewUser("u4"), chat.NewUser("u5")

	server.AddUser(user1)
	server.AddUser(user2)
	server.AddUser(user3)
	server.AddUser(user4)
	server.AddUser(user5)

	channelName := "r"
	server.AddChannel(channelName)
	channel, err := server.GetChannel(channelName)
	if err != nil {
		t.Fatal("Couldn't get the channel just added")
	}

	channel.AddUser("u1")
	channel.AddUser("u2")
	channel.AddUser("u3")
	channel.AddUser("u4")
	channel.AddUser("u5")

	user1.SetChannel(channelName)
	user2.SetChannel(channelName)
	user3.SetChannel(channelName)
	user4.SetChannel(channelName)
	user5.SetChannel(channelName)

	go func(t *testing.T) {
		if _, err := user1.HandleNewInput(server, `/list`); err != nil {
			t.Fatalf("Failed executing message. Error %s", err)
		}
	}(t)

	incoming := user1.GetOutgoing()
	if strings.Contains(incoming, `u1`) {
		t.Errorf("Should not list the executing user in list command. %s should not contain %s", incoming, "u1")
	}
	if !strings.Contains(incoming, `u2`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u2`, incoming)
	}
	if !strings.Contains(incoming, `u3`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u3`, incoming)
	}
	if !strings.Contains(incoming, `u4`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u4`, incoming)
	}
	if !strings.Contains(incoming, `u5`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u5`, incoming)
	}

	// Test list command with another user
	go func(t *testing.T) {
		if _, err := user2.HandleNewInput(server, `/list`); err != nil {
			t.Fatalf("Failed executing message. Error %s", err)
		}
	}(t)

	incoming = user2.GetOutgoing()
	if strings.Contains(incoming, `u2`) {
		t.Errorf("Should not list the executing user in list command. %s should not contain %s", incoming, "u2")
	}
	if !strings.Contains(incoming, `u1`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u2`, incoming)
	}
	if !strings.Contains(incoming, `u3`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u3`, incoming)
	}
	if !strings.Contains(incoming, `u4`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u4`, incoming)
	}
	if !strings.Contains(incoming, `u5`) {
		t.Errorf("Did not get other user in list output. Expected %s got %s", `u5`, incoming)
	}

}
