package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
)

var chatServerInstance *chat.Server

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = new(websocket.Upgrader)
	chatConnection := NewChatConnection()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	chatConnection.Connection = conn
	go listen(chatConnection)
	chatServerInstance.ReceiveConnection(chatConnection)
}

// Start starts chat server
func Start(chatServerParam *chat.Server) {
	chatServerInstance = chatServerParam
}
