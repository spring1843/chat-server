package main

import (
	"log"
	"os"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/rest"
	"github.com/spring1843/chat-server/telnet"
	"github.com/spring1843/chat-server/websocket"
	"github.com/spring1843/pomain/src/shared/errs"
)

func bootstrap(config config.Config) {
	chatServer := chat.NewServer()

	if err := setLogFile(config.LogFile, chatServer); err != nil {
		log.Printf("Error - opening log file %s : %v", config.LogFile, err)
	}

	chatServer.Listen()
	log.Printf("Info - Chat sServer started")

	if err := startTelnet(config, chatServer); err != nil {
		log.Fatalf("Could not start telnet server. Error %s", err)
	}

	startWebSocket(config, chatServer)

	startRESTFulAPI(config, chatServer)
}

func setLogFile(logFile string, chatServer *chat.Server) error {
	if logFile == `` {
		return errs.New("Logfile can not be empty")
	}

	logger, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	chatServer.SetLogFile(logger)
	log.Printf("Info - Log files written to %s", logFile)
	return nil
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
