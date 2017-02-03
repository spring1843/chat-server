package telnet_test

import (
	"bufio"
	"net"
	"testing"

	"strings"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/telnet"
)

func TestCanStartTelnetAndConnectToIt(t *testing.T) {
	config := config.Config{
		TelnetAddress: `0.0.0.0:4002`,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	err := telnet.Start(chatServer, config)
	if err != nil {
		t.Errorf("Could not start telnet server. Error %s", err)
	}

	conn, err := net.Dial("tcp", config.TelnetAddress)
	defer conn.Close()
	if err != nil {
		t.Errorf("Could not connect to the telnet server. Error %s", err)
	}

	welcomeMessage, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("Could not read from the telnet server. Error %s", err)
	}

	if !strings.Contains(welcomeMessage, `Welcome`) {
		t.Errorf("Could not receive welcome message. Message %s", welcomeMessage)
	}
}

func TestOutputErrorWhenCantStart(t *testing.T) {
	config := config.Config{
		TelnetAddress: `0.0.0.0:-1`,
	}

	chatServer := chat.NewServer()

	err := telnet.Start(chatServer, config)
	if err == nil {
		t.Errorf("Server started on an invalid port. Error %s", err)
	}
}
