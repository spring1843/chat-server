package rest

import (
	"github.com/spring1843/chat-server/libs/go-restful"
	"github.com/spring1843/chat-server/libs/go-restful-swagger12"
)

// ConfigureSwagger configures the swagger documentation for all endpoints in the container
func ConfigureSwagger(apiDocPath string, container *restful.Container) {
	if apiDocPath == "" {
		return
	}
	config := swagger.Config{
		WebServices:    container.RegisteredWebServices(),
		WebServicesUrl: ``,
		ApiPath:        apiDocPath,
	}
	swagger.RegisterSwaggerService(config, container)
}
