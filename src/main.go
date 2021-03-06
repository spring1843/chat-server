package main

import (
	"flag"
	"time"

	"os"
	"path/filepath"

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
	config = setCWD(config)
	checkStaticDirExists(config)

	// Start all services e.g. Telnet, WebSocket, REST
	bootstrap.NewBootstrap(config)

	// Never end
	neverDie()
}

func checkStaticDirExists(config config.Config) {
	absolutePath, err := filepath.Abs(filepath.Join(config.CWD, config.StaticWeb))
	if err != nil {
		logs.Fatalf("Error finding absolutepath of %q + %q", config.CWD, config.StaticWeb)
	}
	_, err = os.Stat(absolutePath)
	if os.IsNotExist(err) {
		logs.Fatalf("Directory for StaticWeb defined in config does not exist. CWD %s Absolute Path %s", config.CWD, absolutePath)
		return
	}
}

func setCWD(config config.Config) config.Config {
	if config.CWD == "" {
		executable, err := os.Executable()
		if err != nil {
			logs.FatalIfErrf(err, "Error finding out CWD, current working directory")
		}
		config.CWD += filepath.Dir(executable)
	}
	return config
}

func neverDie() {
	for true {
		time.Sleep(24 * time.Hour)
	}
}
