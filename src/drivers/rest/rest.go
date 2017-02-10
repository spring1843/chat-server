package rest

import (
	"net/http"

	"github.com/spring1843/chat-server/libs/go-restful"
	"github.com/spring1843/chat-server/libs/go-restful-swagger12"
	"github.com/spring1843/chat-server/src/chat"
)

// messageEndpoint an instance of the chat.Server
type messageEndpoint struct {
	ChatServer *chat.Server
}

// LogFilePath path to API log file
var LogFilePath string

func registerAllEndpoints(chatServer *chat.Server, container *restful.Container) {
	messageResource := new(messageEndpoint)
	messageResource.ChatServer = chatServer
	messageResource.Register(container)

	statusResource := new(statusEndpoint)
	statusResource.Register(container)
}

func configureSwagger(wsContainer *restful.Container) swagger.Config {
	return swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(),
		WebServicesUrl: ``,
		ApiPath:        "/api/docs.json",
	}
}

// GetHandler returns a handler that includes all API endpoins
func GetHandler(chatServer *chat.Server) http.Handler {
	handler := restful.NewContainer()
	registerAllEndpoints(chatServer, handler)
	swagger.RegisterSwaggerService(configureSwagger(handler), handler)
	return handler
}
