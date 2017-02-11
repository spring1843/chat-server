package longtests

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	gorilla "github.com/spring1843/chat-server/libs/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/logs"
)

func connectUser(nickname string, wsPath string, config config.Config, i int) *gorilla.Conn {
	url := url.URL{Scheme: "wss", Host: config.WebAddress, Path: wsPath}
	var err error
	var conn *gorilla.Conn
	for ii := 0; ii < dialAttempts; ii++ {
		conn, _, err = gorilla.DefaultDialer.Dial(url.String(), nil)
		logs.ErrIfErrf(err, "Dial attempt %d failed for user%d", ii, i)
	}
	if err != nil {
		logs.Fatalf("user%d error, Websocket couldn't dial after %d attempts %q error: %s", i, dialAttempts, url.String(), err.Error())
	}

	message := readAndIgnoreOtherUserJoinMessages(conn, i)
	if !strings.Contains(message, "Welcome") {
		logs.Fatalf("user%d error, Could not receive welcome message. In %s", i, message)
	}

	if err := conn.WriteMessage(1, []byte(nickname)); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	message = readAndIgnoreOtherUserJoinMessages(conn, i)
	expect := "Welcome " + nickname
	if !strings.Contains(message, expect) {
		logs.Fatalf("user%d error, Could not set user %s, expected 'Thanks User1' got %s", i, nickname, expect)
	}

	return conn
}

func joinChannel(conn *gorilla.Conn, i int) {
	if err := conn.WriteMessage(1, []byte("/join #r")); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	message := readAndIgnoreOtherUserJoinMessages(conn, i)
	expect := fmt.Sprintf("%02d", plugins.UserOutPutTUserCommandOutput)
	if !strings.Contains(message, expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected 'you are now in' (%q) got %q", i, expect, message)
	}

	message = readAndIgnoreOtherUserJoinMessages(conn, i)
	expect = fmt.Sprintf("%02d", plugins.UserOutPutTypeFERunFunction)
	if !strings.Contains(message, expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected 'setChannel' (%q) got %q", i, expect, message)
	}

	message = readAndIgnoreOtherUserJoinMessages(conn, i)
	expect = fmt.Sprintf("%02d", plugins.UserOutPutTUserTraffic)
	if !strings.Contains(message, expect) {
		logs.Fatalf("user%d error, Could not join channel #r. Expected 'You are the first or there are other users' %q got %q", i, expect, message)
	}
}

func readAndIgnoreOtherUserJoinMessages(conn *gorilla.Conn, i int) string {
	_, message, err := conn.ReadMessage()
	if err != nil {
		logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
	}

	// Ignore user traffic messages the user will receive while in the channel.
	inoreUserTraffic := "just joined channel"
	for strings.Contains(string(message), inoreUserTraffic) {
		_, message, err = conn.ReadMessage()
		if err != nil {
			logs.Fatalf("user%d error, Error while reading connection. Error %s", i, err.Error())
		}
	}

	return string(message)
}

func disconnectUser(conn *gorilla.Conn, chatServer *chat.Server, i int) {
	if err := conn.WriteMessage(1, []byte(`/quit`)); err != nil {
		logs.Fatalf("user%d error, Error writing to connection. Error %s", i, err)
	}

	if chatServer.IsUserConnected("User1") {
		log.Fatal("User is still connected to server after quiting")
	}
}
