package telnet_test

import (
	"bufio"
	"net"
	"strconv"
	"testing"

	"strings"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/drivers/telnet"
)

func TestCanStartTelnetAndConnectToIt(t *testing.T) {
	config := config.Config{
		IP:         `0.0.0.0`,
		TelnetPort: 4000,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	err := telnet.Start(chatServer, config)
	if err != nil {
		t.Errorf("Could not start telnet server")
	}

	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(config.TelnetPort))
	defer conn.Close()
	if err != nil {
		t.Errorf("Could not connect to the telnet server")
	}

	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Errorf("Could not read from the telnet server")
	}

	if strings.Contains(status, `Welcome`) != true {
		t.Errorf("Could not receive welcome message")
	}
}

func TestOutputErrorWhenCantStart(t *testing.T) {
	config := config.Config{
		IP:         `0.0.0.0`,
		TelnetPort: -1,
	}

	chatServer := chat.NewServer()

	err := telnet.Start(chatServer, config)
	if err == nil {
		t.Errorf("Server started on an invalid port")
	}
}
