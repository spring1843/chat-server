package bootstrap

import (
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/shared/logs"
)

var chatServer *chat.Server

// NewBootstrap bootstraps chat server and starts all the drivers
func NewBootstrap(config config.Config) {
	chatServer = chat.NewServer()
	chatServer.Listen()
	logs.Info("Chat Server started")

	if config.TelnetAddress != "" {
		logs.FatalIfErrf(startTelnet(config), "Could not start telnet server.")
	} else {
		logs.Warnf("TelnetAddress is empty, not running Telnet Driver")
	}

	if config.WebAddress != "" {
		startWeb(config)
	} else {
		logs.Warnf("WebAddress is empty, not running Web Drivers")
	}
}

// GetChatServer returns thr running instance of chat server
func GetChatServer() *chat.Server {
	return chatServer
}

func startTelnet(config config.Config) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err
	}
	logs.Info("Telnet server started")
	return nil
}
