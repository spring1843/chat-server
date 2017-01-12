package rest

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
)

// messageEndpoint holds an instance of the chat.Server
type messageEndpoint struct {
	ChatServer *chat.Server
}

// LogFilePath is the path to API log file
var LogFilePath string

// Register all rest routes
func registerRoutes(chatServer *chat.Server, container *restful.Container) {
	messageResource := new(messageEndpoint)
	messageResource.ChatServer = chatServer
	messageResource.Register(container)

	statusResource := new(statusEndpoint)
	statusResource.Register(container)
}

// Configure swagger
func configureSwagger(wsContainer *restful.Container) swagger.Config {
	return swagger.Config{
		WebServices:     wsContainer.RegisteredWebServices(),
		WebServicesUrl:  ``,
		ApiPath:         "/docs.json",
		SwaggerPath:     "/docs/",
		SwaggerFilePath: "rest/docs-web-ui",
	}
}

// Start the rest server and configures it
func Start(chatServer *chat.Server, config config.Config) {

	LogFilePath = config.LogFile

	wsContainer := restful.NewContainer()
	registerRoutes(chatServer, wsContainer)
	swagger.RegisterSwaggerService(configureSwagger(wsContainer), wsContainer)

	chatServer.LogPrintf("info \t Rest server listening=%s:%d\nBrowse http://%s:%d/docs/ for RESTful endpoint docs", config.IP, config.RestPort, config.IP, config.RestPort)

	server := &http.Server{Addr: ":" + strconv.Itoa(config.RestPort), Handler: wsContainer}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}
