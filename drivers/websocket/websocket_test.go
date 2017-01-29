package websocket_test

import (
	"net/url"
	"strconv"
	"strings"
	"testing"

	gorilla "github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/drivers/websocket"
)

func TestCantStartAndConnect(t *testing.T) {
	config := config.Config{
		IP:            `0.0.0.0`,
		WebsocketPort: 4008,
	}

	chatServer := chat.NewServer()
	chatServer.Listen()

	websocket.Start(chatServer, config)

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:" + strconv.Itoa(config.WebsocketPort), Path: "/ws"}

	conn, _, err := gorilla.DefaultDialer.Dial(u.String(), nil)
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

	if err := conn.WriteMessage(1, []byte(`User1`)); err != nil {
		t.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("Error while reading connection. Error %s", err.Error())
	}
	if !strings.Contains(string(message), "Thanks User1") {
		t.Fatalf("Could not set user nickname, expected 'Thanks User1' got %s", string(message))
	}

	if err := conn.WriteMessage(1, []byte(`/quit`)); err != nil {
		t.Fatalf("Error writing to connection. Error %s", err)
	}

	_, message, err = conn.ReadMessage()
	if !strings.Contains(string(message), "Good Bye") {
		t.Fatalf("Could not quit from server. Expected 'Good Bye' got %s", string(message))
	}

	if chatServer.IsUserConnected("User1") {
		t.Fatal("User is still connected to server after quiting")
	}
}
