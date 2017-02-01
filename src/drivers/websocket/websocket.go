package websocket

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const templatePath = "drivers/websocket/client/index.html"

var (
	chatServerInstance *chat.Server
	clientTemplate     *template.Template
)

func readTemplate(templatePath string) {
	cwd, err := os.Getwd()
	logs.FatalIfErrf(err, "Failed getting CWD to render client for web socket")


	templateFile := filepath.Join(cwd, templatePath)
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		logs.FatalIfErrf(err, "Failed getting CWD to render client for web socket")
	}

	clientTemplate, err = template.ParseFiles(templateFile)
	logs.FatalIfErrf(err, "Failed reading from %s", templateFile)

	logs.Infof("Read template file for WebSocket client from %s", templatePath)
}

func serveClient(w http.ResponseWriter, r *http.Request) {
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
	chatServerInstance.ReceiveConnection(chatConnection)
}

// Start starts chat server
func Start(chatServerParam *chat.Server, config config.Config) {
	readTemplate(templatePath)
	chatServerInstance = chatServerParam
	http.HandleFunc("/client", serveClient)
	http.HandleFunc("/ws", serveWebSocket)

	go func() {
		err := http.ListenAndServe(config.IP+`:`+strconv.Itoa(config.WebsocketPort), nil)
		if err != nil {
			log.Fatalf("Could not open websocket connection. Error %s", err)
		}
	}()

	logs.Infof("Websocket server listening=%s:%d", config.IP, config.WebsocketPort)
	logs.Infof("Browse http://%s:%d/client/ for Websocket client", config.IP, config.WebsocketPort)
}
