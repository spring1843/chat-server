package main

import (
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/rest"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

func bootstrap(config config.Config) {
	chatServer := chat.NewServer()
	chatServer.Listen()
	logs.Info("Chat Server started")

	logs.FatalIfErrf(startTelnet(config, chatServer), "Could not start telnet server.")

	startWebSocket(config, chatServer)

	startRESTFulAPI(config, chatServer)
}

func startTelnet(config config.Config, chatServer *chat.Server) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err

	}
	logs.Info("Telnet server started")
	return nil
}

func startWebSocket(config config.Config, chatServer *chat.Server) {
	websocket.Start(chatServer, config)
	logs.Info("WebSocket server started")
}

// startRESTFulAPI starts the restful API
func startRESTFulAPI(config config.Config, chatServer *chat.Server) {
	server := rest.NewRESTfulAPI(config, chatServer)
	go func() {
		logs.FatalIfErrf(server.ListenAndServe(), "Could not start Rest server.")
	}()
	logs.Infof("RESTful API started")
}
