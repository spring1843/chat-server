package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
)

var chatServerInstance *chat.Server

// Handler is a http handler function that implements WebSocket
func Handler(w http.ResponseWriter, r *http.Request) {
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

// SetWebSocket sets the chat server instance
func SetWebSocket(chatServerParam *chat.Server) {
	chatServerInstance = chatServerParam
}
