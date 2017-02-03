package websocket_test

import (
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

func TestCantStartAndConnect(t *testing.T) {
	config := config.Config{
		WebAddress: "127.0.0.1:4008",
	}

	chatServer := chat.NewServer()
	chatServer.Listen()
	websocket.Start(chatServer)

	go func() {
		err := http.ListenAndServe(config.WebAddress, nil)
		if err != nil {
			log.Fatalf("Could not open websocket connection. Error %s", err)
		}
	}()
	u := url.URL{Scheme: "ws", Host: config.WebAddress, Path: "/ws"}

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

	expect := "Welcome User1"
	if !strings.Contains(string(message), expect) {
		t.Fatalf("Could not set user nickname, expected 'Thanks User1' got %s", expect)
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
