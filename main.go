package main

import (
	"flag"
	"log"

	"github.com/spring1843/chat-server/config"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "path to .json config file")
	flag.Parse()

	if configFile == "" {
		log.Fatal("-config flag is required")
	}
	config := config.FromFile(configFile)

	// Start all services e.g. Telnet, WebSocket, REST
	bootstrap(config)

	// Never end
	for true {
	}
}
