package websocket_test

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	gorilla "github.com/spring1843/chat-server/libs/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const dialAttempts = 3

func TestCantStartAndConnect(t *testing.T) {
	if os.Getenv("SKIP_NETWORK") == "1" {
		t.Skipf("Skipping test SKIP_NETWORK set to %q", os.Getenv("SKIPNETWORK"))
	}
	config := config.Config{
		WebAddress: "127.0.0.1:4003",
	}

	chatServer := chat.NewServer()
	chatServer.Listen()
	websocket.SetWebSocket(chatServer)

	http.HandleFunc("/ws1", websocket.Handler)

	go func() {
		if err := http.ListenAndServe(config.WebAddress, nil); err != nil {
			t.Fatalf("Failed listening to Websocet on %s. Error: %s", config.WebAddress, err)
		}
	}()
	u := url.URL{Scheme: "ws", Host: config.WebAddress, Path: "/ws1"}

	var err error
	var conn *gorilla.Conn
	for i := 0; i < dialAttempts; i++ {
		conn, _, err = gorilla.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logs.ErrIfErrf(err, "Dial attempt %d failed", i)
		}
	}
	if err != nil {
		t.Fatalf("Websocket dial attempt failed after %d times. error: %s", dialAttempts, err.Error())
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
