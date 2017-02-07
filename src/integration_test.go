package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"

	gorilla "github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

var (
	integrationTryouts = 100
	doneWithAllUsers   chan bool
)

// Tests that the server can be started with config.json configs
// And many users can connect to it using WebSocket, join a channel, chat and then disconnect
func TestManyUsers(t *testing.T) {
	config := config.FromFile("./config.json")
	config.WebAddress += "3"
	config.TelnetAddress += "3"

	bootstrap(config)

	i := 0
	for i < integrationTryouts {
		nickName := fmt.Sprintf("user%d", i)
		go connectAndDisconnect(nickName, "/ws", config, chatServer, i)
		i++
	}
	<-doneWithAllUsers
	close(doneWithAllUsers)
}

func connectUser(nickname string, wsPath string, config config.Config, i int) *gorilla.Conn {
	url := url.URL{Scheme: "wss", Host: config.WebAddress, Path: wsPath}

	conn, _, err := gorilla.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		logs.Fatalf("user%d error, Websocket couldn't dial %q error: %s", i, url.String(), err.Error())
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection %s", i, err.Error())
	}

	if !strings.Contains(string(message), "Welcome") {
		logs.Fatalf("user%d error, Could not receive welcome message. In %s", i, message)
	}

	if err := conn.WriteMessage(1, []byte(nickname)); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
	}

	expect := "Welcome " + nickname
	if !strings.Contains(string(message), expect) {
		logs.Fatalf("user%d error, Could not set user %s, expected 'Thanks User1' got %s", i, nickname, expect)
	}

	return conn
}

func joinChannel(conn *gorilla.Conn, i int) {
	if err := conn.WriteMessage(1, []byte("/join #r")); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
	}
	expect := "05"
	if !strings.Contains(string(message), expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected %q got %q", i, expect, message)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
	}
	expect = "06"
	if !strings.Contains(string(message), expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected %q got %q", i, expect, message)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
	}
	expect = "00"
	if !strings.Contains(string(message), expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected %q got %q", i, expect, message)
	}
}

func disconnectUser(conn *gorilla.Conn, chatServer *chat.Server, i int) {
	if err := conn.WriteMessage(1, []byte(`/quit`)); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	if _, _, err := conn.ReadMessage(); err != nil {
		logs.Fatalf("user%d error, Failed reading from WebSocket connection. Error %s", i, err)
	}

	if chatServer.IsUserConnected("User1") {
		log.Fatal("User is still connected to server after quiting")
	}
}

func connectAndDisconnect(nickname string, wsPath string, config config.Config, chatServer *chat.Server, i int) {
	conn := connectUser(nickname, wsPath, config, i)
	joinChannel(conn, i)
	disconnectUser(conn, chatServer, i)
	if i == integrationTryouts {
		doneWithAllUsers <- true
	}
}
