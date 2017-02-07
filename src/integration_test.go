package main

import (
	"log"
	"net/url"
	"strings"
	"testing"

	"fmt"
	"sync"
	"time"

	gorilla "github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/logs"
)

var (
	curUserCount = 0
	userCount    = 100
	counterLock  = new(sync.Mutex)

	doneWithAllUsers chan bool
)

// Tests that the server can be started with config.json configs
// And many users can connect to it using WebSocket, join a channel, chat and then disconnect
func TestManyUsers(t *testing.T) {
	config := config.FromFile("./config.json")
	config.WebAddress += "3"
	config.TelnetAddress = ""

	bootstrap(config)

	doneWithAllUsers = make(chan bool, 1)

	i := 0
	for i < userCount {
		nickName := fmt.Sprintf("user%d", i)
		go connectAndDisconnect(nickName, "/ws", config, chatServer, i)
		i++
	}

	select {
	case <-time.After(time.Second * 15):
		t.Fatalf("Didn't finish after 5 seconds")
	case done := <-doneWithAllUsers:
		if done {
			t.Logf("Done!")
		}
	}
}

func connectUser(nickname string, wsPath string, config config.Config, i int) *gorilla.Conn {
	url := url.URL{Scheme: "wss", Host: config.WebAddress, Path: wsPath}

	conn, _, err := gorilla.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		logs.Fatalf("user%d error, Websocket couldn't dial %q error: %s", i, url.String(), err.Error())
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

func connectAndDisconnect(nickname string, wsPath string, config config.Config, chatServer *chat.Server, i int) {
	conn := connectUser(nickname, wsPath, config, i)
	joinChannel(conn, i)
	disconnectUser(conn, chatServer, i)

	counterLock.Lock()
	defer counterLock.Unlock()

	curUserCount++
	if curUserCount == userCount {
		doneWithAllUsers <- true
	}
}
