package main

import (
	"net/http"

	"crypto/tls"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/rest"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const staticWebContentDir = "../bin/web"

var chatServer *chat.Server

func bootstrap(config config.Config) {
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

func startTelnet(config config.Config) error {
	err := telnet.Start(chatServer, config)
	if err != nil {
		return err
	}
	logs.Info("Telnet server started")
	return nil
}

func startWeb(config config.Config) {
	go func() {
		srv := getTLSServer(getMultiplexer(), config.WebAddress)
		logs.Infof("Serving static files, Rest, WebSocket on http:/%s/", config.WebAddress)
		logs.FatalIfErrf(srv.ListenAndServeTLS("tls.crt", "tls.key"), "Could not start Rest server.")
	}()
}

func getMultiplexer() *http.ServeMux {
	restHandler := rest.GetHandler(chatServer)
	websocket.SetWebSocket(chatServer)
	fs := http.FileServer(http.Dir(staticWebContentDir))

	mux := http.NewServeMux()
	mux.Handle("/api/", restHandler)
	mux.HandleFunc("/ws", websocket.Handler)
	mux.Handle("/", fs)
	return mux
}

func getTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		//CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

func getTLSServer(mux *http.ServeMux, webAddress string) *http.Server {
	return &http.Server{
		Addr:         webAddress,
		Handler:      mux,
		TLSConfig:    getTLSConfig(),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}
