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
)

// Tests that the server can be started with config.json configs
// And many users can connect to it using WebSocket, join a channel, chat and then disconnect
func TestManyUsers(t *testing.T) {
	config := config.FromFile("./config.json")
	bootstrap(config)

	tryouts := 100
	i := 0
	for i < tryouts {
		nickName := fmt.Sprintf("user%d", i)
		go connectAndDisconnect(nickName, "/ws", config, chatServer)
		i++
	}
}

func connectUser(nickname string, wsPath string, config config.Config) *gorilla.Conn {
	url := url.URL{Scheme: "wss", Host: config.WebAddress, Path: wsPath}

	conn, _, err := gorilla.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatalf("Websocket couldn't dial %q error: %s", url.String(), err.Error())
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("Error while reading connection %s", err.Error())
	}

	if !strings.Contains(string(message), "Welcome") {
		log.Fatalf("Could not receive welcome message. In %s", message)
	}

	if err := conn.WriteMessage(1, []byte(nickname)); err != nil {
		log.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		log.Fatalf("Error while reading connection. Error %s", err.Error())
	}

	expect := "Welcome " + nickname
	if !strings.Contains(string(message), expect) {
		log.Fatalf("Could not set user %s, expected 'Thanks User1' got %s", nickname, expect)
	}

	return conn
}

func joinChannel(conn *gorilla.Conn) {
	if err := conn.WriteMessage(1, []byte("/join #r")); err != nil {
		log.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("Error while reading connection. Error %s", err.Error())
	}
	expect := "05"
	if !strings.Contains(string(message), expect) {
		log.Fatalf("Could not join channel #r. Expected %q got %q", expect, message)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		log.Fatalf("Error while reading connection. Error %s", err.Error())
	}
	expect = "06"
	if !strings.Contains(string(message), expect) {
		log.Fatalf("Could not join channel #r. Expected %q got %q", expect, message)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		log.Fatalf("Error while reading connection. Error %s", err.Error())
	}
	expect = "00"
	if !strings.Contains(string(message), expect) {
		log.Fatalf("Could not join channel #r. Expected %q got %q", expect, message)
	}
}

func disconnectUser(conn *gorilla.Conn, chatServer *chat.Server) {
	if err := conn.WriteMessage(1, []byte(`/quit`)); err != nil {
		log.Fatalf("Error writing to connection. Error %s", err)
	}

	if _, _, err := conn.ReadMessage(); err != nil {
		log.Fatalf("Failed reading from WebSocket connection. Error %s", err)
	}

	if chatServer.IsUserConnected("User1") {
		log.Fatal("User is still connected to server after quiting")
	}
}

func connectAndDisconnect(nickname string, wsPath string, config config.Config, chatServer *chat.Server) {
	conn := connectUser(nickname, wsPath, config)
	joinChannel(conn)
	disconnectUser(conn, chatServer)
}
