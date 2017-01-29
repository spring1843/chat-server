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

func Test_CantStartAndConnect(t *testing.T) {
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
		t.Errorf("Websocket Dial error: %s", err.Error())
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("Error while reading connection %s", err.Error())
	}

	if strings.Contains(string(message), "Welcome") != true {
		t.Error("Could not receive welcome message")
	}
	conn.WriteMessage(1, []byte(`User1`))
	_, message, err = conn.ReadMessage()
	if err != nil {
		t.Errorf("Error while reading connection %s", err.Error())
	}
	if strings.Contains(string(message), "Thanks User1") != true {
		t.Error("Could not set user nickname")
	}

	conn.WriteMessage(1, []byte(`/quit`))
	_, message, err = conn.ReadMessage()

	if err == nil {
		t.Error("Connection didn't close after running quit")
	}
}
