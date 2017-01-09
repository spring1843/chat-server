package rest

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
)

type RestMessageResource struct {
	ChatServer *chat.Server
}

var LogFilePath string

// Register all rest routes
func registerRoutes(chatServer *chat.Server, container *restful.Container) {
	messageResource := new(RestMessageResource)
	messageResource.ChatServer = chatServer
	messageResource.Register(container)

	statusResource := new(StatusResource)
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

// Starts the rest server and configures it
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
