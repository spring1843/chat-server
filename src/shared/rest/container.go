package rest

import "github.com/spring1843/chat-server/libs/go-restful"

// Container interfaces *restful.Container
type Container interface {
	Add(service *restful.WebService) *restful.Container
	RegisteredWebServices() []*restful.WebService
}

// NewHTTPHandler returns a new restful container
func NewHTTPHandler() *restful.Container {
	return restful.NewContainer()
}
