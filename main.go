package main

import (
	"flag"

	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/plugins/logs"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "path to .json config file")
	flag.Parse()

	if configFile == "" {
		logs.Fatalf("-config flag is required")
	}
	config := config.FromFile(configFile)

	// Start all services e.g. Telnet, WebSocket, REST
	bootstrap(config)

	// Never end
	for true {
	}
}
