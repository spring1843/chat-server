package rest

import restful "github.com/emicklei/go-restful"

// Container interfaces *restful.Container, a web service container
type Container interface {
	Add(service *restful.WebService) *restful.Container
	RegisteredWebServices() []*restful.WebService
}

// NewHTTPHandler returns a new restful container
func NewHTTPHandler() *restful.Container {
	return restful.NewContainer()
}
