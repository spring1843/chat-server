package main

import (
	"net/http"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/rest"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const staticWebContentDir = "../bin/web"

func bootstrap(config config.Config) {
	chatServer := chat.NewServer()
	chatServer.Listen()
	logs.Info("Chat Server started")

	if config.TelnetAddress != "" {
		logs.FatalIfErrf(startTelnet(config, chatServer), "Could not start telnet server.")
	} else {
		logs.Warnf("TelnetAddress is empty, not running Telnet Driver")
	}

	if config.WebAddress != "" {
		logs.FatalIfErrf(startWeb(config, chatServer), "Could not start web server.")
	} else {
		logs.Warnf("WebAddress is empty, not running Web Drivers")
	}
}

func startTelnet(config config.Config, chatServer *chat.Server) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err

	}
	logs.Info("Telnet server started")
	return nil
}

func startWeb(config config.Config, chatServer *chat.Server) error {
	restHandler := rest.GetHandler(chatServer)
	websocket.SetWebSocket(chatServer)
	fs := http.FileServer(http.Dir(staticWebContentDir))

	http.Handle("/api/", restHandler)
	http.HandleFunc("/ws", websocket.Handler)
	http.Handle("/", fs)

	go func() {
		logs.Infof("Serving static files, Rest, WebSocket on http:/%s/", config.WebAddress)
		logs.FatalIfErrf(http.ListenAndServeTLS(config.WebAddress, "tls.crt", "tls.key", nil), "Could not start Rest server.")
	}()
	return nil
}
