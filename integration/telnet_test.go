package integration_test

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"testing"

	"strings"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/integration"
	"github.com/spring1843/chat-server/telnet"
)

func TestCanStartTelnetAndConnectToIt(t *testing.T) {
	config := integration.Config{
		IP:         `0.0.0.0`,
		TelnetPort: 4000,
		LogFile:    `/dev/null`,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	testFile, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	chatServer.SetLogFile(testFile)

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
	config := integration.Config{
		IP:         `0.0.0.0`,
		TelnetPort: -1,
		LogFile:    `/dev/null`,
	}

	chatServer := chat.NewServer()

	testFile, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	chatServer.SetLogFile(testFile)

	err := telnet.Start(chatServer, config)
	if err == nil {
		t.Errorf("Server started on an invalid port")
	}
}
