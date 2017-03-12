package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/shared/logs"
)

var (
	chatServerInstance *chat.Server
	upgrader           = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

// TODO validate CORS headers here
func checkOrigin(r *http.Request) bool {
	return true
}

// Handler is a http handler function that implements WebSocket
func Handler(w http.ResponseWriter, r *http.Request) {
	logs.Infof("Call to websocket /wp form %s", r.RemoteAddr)
	chatConnection := NewChatConnection()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logs.ErrIfErrf(err, "Error upgrading websocket connection.")
		return
	}
	chatConnection.Connection = conn
	go chatServerInstance.ReceiveConnection(chatConnection)
	go chatConnection.writePump()
	chatConnection.readPump()
	logs.Infof("End of call to websocket /wp form %s", r.RemoteAddr)
}

// SetWebSocket sets the chat server instance
func SetWebSocket(chatServerParam *chat.Server) {
	chatServerInstance = chatServerParam
}
