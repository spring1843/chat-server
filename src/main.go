package main

import (
	"flag"
	"time"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
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
		time.Sleep(24 * time.Hour)
	}
}
