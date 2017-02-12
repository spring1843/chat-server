package webapi

import (
	"net/http"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/shared/rest"
)

var (
	// LogFilePath path to API log file
	LogFilePath string

	chatServerInstance *chat.Server
)

func registerPaths(container rest.Container) {
	registerMessagePath(container)
	registerStatusPath(container)
}

// NewHandler returns a HTTP handler that includes all RESTfyk API endpoints exposed
func NewHandler(chatServer *chat.Server) http.Handler {
	chatServerInstance = chatServer

	handler := rest.NewHTTPHandler()
	registerPaths(handler)
	rest.ConfigureSwagger("/api/docs.json", handler)
	return handler
}
