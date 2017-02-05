package websocket_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gorilla "github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/websocket"
)

func TestCantStartTwoUsers(t *testing.T) {
	config := config.Config{
		WebAddress: "127.0.0.1:4008",
	}

	chatServer := chat.NewServer()
	chatServer.Listen()
	websocket.SetWebSocket(chatServer)

	http.HandleFunc("/ws1", websocket.Handler)

	go func() {
		if err := http.ListenAndServe(config.WebAddress, nil); err != nil {
			t.Fatalf("Failed listening to WebSocket on %s. Error %s.", config.WebAddress, err)
		}
	}()

	tryouts := 2
	conns := make([]*gorilla.Conn, tryouts, tryouts)
	i := 0
	for i < tryouts {
		nickName := fmt.Sprintf("user%d", i)
		conns[i] = connectUser(nickName, "/ws1", config)
		i++
	}

	if chatServer.ConnectedUsersCount() != tryouts {
		t.Fatalf("Expected user count to be %d after disconnecting users, got %d", tryouts, chatServer.ConnectedUsersCount())
	}

	i = 0
	for i < tryouts {
		disconnectUser(conns[i], chatServer)
		i++
	}

	if chatServer.ConnectedUsersCount() != 0 {
		t.Fatalf("Expected user count to be %d after disconnecting users, got %d", 0, chatServer.ConnectedUsersCount())
	}
}

func TestCantStartAndConnectManyUsers(t *testing.T) {
	config := config.Config{
		WebAddress: "127.0.0.1:4009",
	}

	chatServer := chat.NewServer()
	chatServer.Listen()
	websocket.SetWebSocket(chatServer)

	http.HandleFunc("/ws2", websocket.Handler)

	go func() {
		if err := http.ListenAndServe(config.WebAddress, nil); err != nil {
			t.Fatalf("Failed listening to WebSocket on %s. Error %s.", config.WebAddress, err)
		}
	}()

	tryouts := 100
	i := 0
	for i < tryouts {
		nickName := fmt.Sprintf("user%d", i)
		go connectAndDisconnect(t, nickName, "/ws2", config, chatServer)
		i++
	}
}

func connectUser(nickname string, wsPath string, config config.Config) *gorilla.Conn {
	url := url.URL{Scheme: "ws", Host: config.WebAddress, Path: wsPath}

	conn, _, err := gorilla.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatalf("Websocket Dial error: %s", err.Error())
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

func connectAndDisconnect(t *testing.T, nickname string, wsPath string, config config.Config, chatServer *chat.Server) {
	conn := connectUser(nickname, wsPath, config)
	joinChannel(conn)
	disconnectUser(conn, chatServer)
}
