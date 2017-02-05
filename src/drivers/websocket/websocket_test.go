package websocket_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	gorilla "github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/websocket"
)

func TestCantStartAndConnect(t *testing.T) {
	config := config.Config{
		WebAddress: "127.0.0.1:4008",
	}

	chatServer := chat.NewServer()
	chatServer.Listen()
	websocket.SetWebSocket(chatServer)

	http.HandleFunc("/ws", websocket.Handler)

	go func() {
		if err := http.ListenAndServe(config.WebAddress, nil); err != nil {
			t.Fatalf("Failed listening to Websocet on %s. Error: %s", config.WebAddress, err)
		}
	}()

	tryouts := 10
	i := 0
	for i < tryouts {
		nickName := fmt.Sprintf("user%d", i)
		go connectAndDisconnect(t, nickName, config, chatServer)
		i++
	}
}

func connectUser(t *testing.T, nickname string, config config.Config) *gorilla.Conn {
	url := url.URL{Scheme: "ws", Host: config.WebAddress, Path: "/ws"}

	conn, _, err := gorilla.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		t.Fatalf("Websocket Dial error: %s", err.Error())
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error while reading connection %s", err.Error())
	}

	if !strings.Contains(string(message), "Welcome") {
		t.Error("Could not receive welcome message")
	}

	if err := conn.WriteMessage(1, []byte(nickname)); err != nil {
		t.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error while reading connection. Error %s", err.Error())
	}

	expect := "Welcome " + nickname
	if !strings.Contains(string(message), expect) {
		t.Fatalf("Could not set user %s, expected 'Thanks User1' got %s", nickname, expect)
	}

	return conn
}

func disconnectUser(t *testing.T, conn *gorilla.Conn, chatServer *chat.Server) {
	if err := conn.WriteMessage(1, []byte(`/quit`)); err != nil {
		t.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed reading from WebSocket connection. Error %s", err)
	}
	if !strings.Contains(string(message), "Good Bye") {
		t.Fatalf("Could not quit from server. Expected 'Good Bye' got %s", string(message))
	}

	if chatServer.IsUserConnected("User1") {
		t.Fatal("User is still connected to server after quiting")
	}
}

func connectAndDisconnect(t *testing.T, nickname string, config config.Config, chatServer *chat.Server) {
	conn := connectUser(t, nickname, config)
	disconnectUser(t, conn, chatServer)
}
