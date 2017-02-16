package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/bootstrap"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const usageDoc = `Chat Server
Usage:
        chat-server -config config.json
Flags:
        -config S r equired .json config file, look at config.json for default settings
`

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "path to .json config file")
	flag.Parse()

	if configFile == "" {
		logs.Fatalf(usageDoc)
	}
	config := config.FromFile(configFile)
	setCWD(config)

	// Start all services e.g. Telnet, WebSocket, REST
	bootstrap.NewBootstrap(config)

	// Never end
	neverDie()
}

func setCWD(config config.Config) {
	if config.CWD == "" {
		var err error
		if config.CWD, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			logs.FatalIfErrf(err, "Error finding out CWD, current working directory")
		}
	}
}

func neverDie() {
	for true {
		time.Sleep(24 * time.Hour)
	}
}
