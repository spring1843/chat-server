package longtests

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/bootstrap"
)

const (
	userCount = 100
	timeOutS  = 20
)

var (
	connectedUsersCount = 0
	counterLock         = new(sync.Mutex)

	doneWithAllUsers chan bool
)

// TestManyUsers tests that many users can connect to a chat server using WebSocket connections
// Each user:
// 	tries to establish a WebSocket connection with a number of attempts
// 	gets identified with a nickname
// 	joins a channel
// 	quits the server
func TestManyUsers(t *testing.T) {
	if os.Getenv("LONGTESTS") != "1" {
		t.Skipf("LONGTESTS not set to 1, run make longtest to run this test")
	}
	if os.Getenv("SKIP_NETWORK") == "1" {
		t.Skipf("Skipping test SKIP_NETWORK set to %q", os.Getenv("SKIPNETWORK"))
	}

	config := config.FromFile("../config.json")
	config.WebAddress += "3"
	config.TelnetAddress = ""

	bootstrap.NewBootstrap(config)
	doneWithAllUsers = make(chan bool, 1)

	for i := 0; i < userCount; i++ {
		nickName := fmt.Sprintf("user%d", i)

		// Each user connects and disconnects concurrently
		go connectAndDisconnect(nickName, "/ws", config, bootstrap.GetChatServer(), i)
	}

	select {
	case <-time.After(timeOutS * time.Second):
		t.Fatalf("Timeout, didn't finish after %d seconds", timeOutS)
	case done := <-doneWithAllUsers:
		if done {
			t.Logf("Done!")
		}
	}
}

func connectAndDisconnect(nickname string, wsPath string, config config.Config, chatServer *chat.Server, i int) {
	conn := connectUser(nickname, wsPath, config, i)
	joinChannel(conn, i)
	disconnectUser(conn, chatServer, i)

	counterLock.Lock()
	defer counterLock.Unlock()

	connectedUsersCount++
	if connectedUsersCount == userCount {
		doneWithAllUsers <- true
	}
}
