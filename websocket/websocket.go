package websocket

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
)

var chatServer *chat.Server

func serveClient(w http.ResponseWriter, r *http.Request) {
	var cwd, _ = os.Getwd()
	var clientTemplate = template.Must(template.ParseFiles(filepath.Join(cwd, "websocket/client/index.html")))
	if r.URL.Path != "/client" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	clientTemplate.Execute(w, r.Host)
}

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	var upgrader = new(websocket.Upgrader)
	chatConnection := NewChatConnection()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	chatConnection.Connection = conn
	go listen(chatConnection)
	chatServer.Connection <- chatConnection
}

// Start starts chat server
func Start(chatServer *chat.Server, config config.Config) error {
	chatServer = chatServer
	http.HandleFunc("/client", serveClient)
	http.HandleFunc("/ws", serveWebSocket)

	go func() {
		err := http.ListenAndServe(config.IP+`:`+strconv.Itoa(config.WebsocketPort), nil)
		if err != nil {
			panic("Could not open websocket connection, address already in use?")
		}

	}()
	chatServer.LogPrintf("info \t Websocket server listening=%s:%d\nBrowse http://%s:%d/client/ for Websocket client", config.IP, config.WebsocketPort, config.IP, config.WebsocketPort)

	return nil
}
