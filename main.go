package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/rest"
	"github.com/spring1843/chat-server/telnet"
	"github.com/spring1843/chat-server/websocket"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "path to .json config file")
	flag.Parse()

	if configFile == "" {
		log.Fatal("-config flag is required")
	}
	config := config.FromFile(configFile)

	chatServer := configureNewChatServer(config)
	startTelnet(config, chatServer)
	startWebSocket(config, chatServer)
	startRest(config, chatServer)

	neverEnd()
}

// configureNewChatServer configures a new chat server
func configureNewChatServer(config config.Config) *chat.Server {
	chatServer := chat.NewServer()
	chatServer.Listen()
	setLogFile(config.LogFile, chatServer)
	return chatServer
}

// startTelnet starts the telnet server
func startTelnet(config config.Config, chatServer *chat.Server) {
	err := telnet.Start(chatServer, config)
	if err != nil {
		log.Printf("Could not start telnet server please check the logs for more info")
		panic(err)
	}
	log.Printf("Info - Telnet Server Started")
}

// startRest starts the rest server
func startRest(config config.Config, chatServer *chat.Server) {
	rest.Start(chatServer, config)
	log.Printf("Info - Web Server Started")
}

// startWebSocket starts the Websocket server
func startWebSocket(config config.Config, chatServer *chat.Server) {
	log.Printf("Info - Websocket Server Started")
	err := websocket.Start(chatServer, config)
	if err != nil {
		log.Printf("Could not start websocket server please check the logs for more info")
		panic(err)
	}
}

// setLogFile creates a log file to be given to a server
func setLogFile(logFile string, chatServer *chat.Server) {
	if logFile == `` {
		return
	}

	logger, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error - opening log file %s : %v", logFile, err)
		panic(err)
	}
	chatServer.SetLogFile(logger)
	log.Printf("Info - Log files written to %s", logFile)
}

// neverEnd never end execution, avoids termination while the server is concurrently running
func neverEnd() {
	log.Printf("Waiting for TCP connections...")
	for {
		time.Sleep(1000 * time.Second)
	}
}
