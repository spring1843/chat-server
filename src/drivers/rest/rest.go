package rest

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/logs"
)

const templatePath = "drivers/rest/docs-web-ui"

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
		WebServices:     wsContainer.RegisteredWebServices(),
		WebServicesUrl:  ``,
		ApiPath:         "/docs.json",
		SwaggerPath:     "/docs/",
		SwaggerFilePath: templatePath,
	}
}

// NewRESTfulAPI the rest server and configures it
func NewRESTfulAPI(config config.Config, chatServer *chat.Server) *http.Server {
	wsContainer := restful.NewContainer()
	registerAllEndpoints(chatServer, wsContainer)
	swagger.RegisterSwaggerService(configureSwagger(wsContainer), wsContainer)

	logs.Infof("Rest server listening=%s:%d", config.IP, config.RestPort)
	logs.Infof("Browse http://%s:%d/docs/ for RESTful endpoint docs", config.IP, config.RestPort)

	return &http.Server{Addr: ":" + strconv.Itoa(config.RestPort), Handler: wsContainer}
}
