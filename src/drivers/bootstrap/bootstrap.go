package bootstrap

import (
	"crypto/tls"
	"net/http"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/telnet"
	"github.com/spring1843/chat-server/src/drivers/webapi"
	"github.com/spring1843/chat-server/src/drivers/websocket"
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

func startWeb(config config.Config) {
	srv := getTLSServer(getMultiplexer(config), config.WebAddress)
	go func() {
		logs.Infof("Serving static files, Rest, WebSocket on http:/%s/", config.WebAddress)
		logs.FatalIfErrf(srv.ListenAndServeTLS("tls.crt", "tls.key"), "Could not start Rest server. Error %s")
	}()
}

func getMultiplexer(config config.Config) *http.ServeMux {
	restHandler := webapi.NewHandler(chatServer)
	websocket.SetWebSocket(chatServer)

	mux := http.NewServeMux()
	mux.Handle("/api/", restHandler)
	mux.HandleFunc("/ws", websocket.Handler)
	serveStaticWeb(mux, config)
	return mux
}

func serveStaticWeb(mux *http.ServeMux, config config.Config) {
	if config.StaticWeb == "" {
		logs.Infof("Not serving static web files")
		return
	}
	logs.Infof("Serving static web files from %s", config.StaticWeb)
	fs := http.FileServer(http.Dir(config.StaticWeb))
	mux.Handle("/", fs)
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
