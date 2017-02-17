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
func NewHandler(chatServer *chat.Server, apiDocPath string) http.Handler {
	chatServerInstance = chatServer

	handler := rest.NewHTTPHandler()
	registerPaths(handler)
	if apiDocPath != "" {
		rest.ConfigureSwagger(apiDocPath, handler)
	}
	return handler
}
