package main

import (
	"net/http"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/rest"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const staticWebContentDir = "../bin/web"

func bootstrap(config config.Config) {
	chatServer := chat.NewServer()
	chatServer.Listen()
	logs.Info("Chat Server started")

	logs.FatalIfErrf(startTelnet(config, chatServer), "Could not start telnet server.")

	restHandler := rest.GetHandler(chatServer)
	websocket.Start(chatServer)
	fs := http.FileServer(http.Dir(staticWebContentDir))

	http.Handle("/api/", restHandler)
	http.HandleFunc("/ws", websocket.WebSocketHandler)
	http.Handle("/", fs)

	go func() {
		logs.FatalIfErrf(http.ListenAndServe(config.WebAddress, nil), "Could not start Rest server.")
	}()
}

func startTelnet(config config.Config, chatServer *chat.Server) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err

	}
	logs.Info("Telnet server started")
	return nil
}
