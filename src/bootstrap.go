package main

import (
	"log"

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
	logs.Infof("Info - Chat sServer started")

	if err := startTelnet(config, chatServer); err != nil {
		log.Fatalf("Could not start telnet server. Error %s", err)
	}

	startWebSocket(config, chatServer)

	startRESTFulAPI(config, chatServer)
}

func startTelnet(config config.Config, chatServer *chat.Server) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err

	}
	log.Printf("Info - Telnet server started")
	return nil
}

func startWebSocket(config config.Config, chatServer *chat.Server) {
	websocket.Start(chatServer, config)
	log.Printf("Info - WebSocket server started")
}

// startRESTFulAPI starts the restful API
func startRESTFulAPI(config config.Config, chatServer *chat.Server) {
	server := rest.NewRESTfulAPI(config, chatServer)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed starting RESTFul API, Error %s", err)
		}
	}()
	log.Printf("Info - RESTful API started")
}
