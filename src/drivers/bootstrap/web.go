package bootstrap

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/drivers/webapi"
	"github.com/spring1843/chat-server/src/drivers/websocket"
	"github.com/spring1843/chat-server/src/shared/logs"
)

func startWeb(config config.Config) {
	srv := getTLSServer(getMultiplexer(config), config.WebAddress)
	go func() {
		var err error

		switch config.HTTPS {
		case false:
			logs.Infof("HTTPS disabled in config")
			logs.Infof("Serving static files, Rest, WebSocket on http:/%s/", config.WebAddress)
			err = srv.ListenAndServe()
		default:
			absolutePathCert, err := filepath.Abs(filepath.Join(config.CWD, config.TLSCert))
			if err != nil {
				logs.Fatalf("Error finding absolute path of TLS cert %s%s", config.CWD, config.TLSCert)
			}
			if _, err = os.Stat(absolutePathCert); os.IsNotExist(err) {
				logs.Fatalf("TLS cert file path defined in config does not exist. CWD %s Absolute Path %s", config.CWD, absolutePathCert)
				return
			}

			absolutePathKey, err := filepath.Abs(filepath.Join(config.CWD, config.TLSKey))
			if err != nil {
				logs.Fatalf("Error finding absolute path of TLS Key %s%s", config.CWD, config.TLSKey)
			}
			_, err = os.Stat(absolutePathKey)
			if os.IsNotExist(err) {
				logs.Fatalf("TLS key file path defined in config does not exist. CWD %s Absolute Path %s", config.CWD, absolutePathKey)
				return
			}
			logs.Infof("Serving static files, Rest, WebSocket on https:/%s/", config.WebAddress)
			err = srv.ListenAndServeTLS("tls.crt", "tls.key")
		}
		logs.FatalIfErrf(err, "Could not start Rest server. Error %s", err)
	}()
}

func getMultiplexer(config config.Config) *http.ServeMux {
	if config.APIDocPath == "" {
		logs.Infof("Not serving API Docs JSON endpoit because APIDocPath is empty in config")
	}
	restHandler := webapi.NewHandler(chatServer, config.APIDocPath)
	websocket.SetWebSocket(chatServer)

	mux := http.NewServeMux()
	mux.Handle("/api/", restHandler)
	mux.HandleFunc("/ws", websocket.Handler)
	serveStaticWeb(mux, config)
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
