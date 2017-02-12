package rest

import "github.com/spring1843/chat-server/libs/go-restful"

type (
	// Container interfaces *restful.Container
	Container interface {
		Add(service *restful.WebService) *restful.Container
		RegisteredWebServices() []*restful.WebService
	}
	// EndpointHandlerParams is a value passed to every function that is supposed to handle RESTful calls
	EndpointHandlerParams struct {
		Req  *restful.Request
		Resp *restful.Response
	}
	// EndpointFunction is a function that intakes EndpointHandlerParams to respond to a RESTful call
	EndpointFunction func(*EndpointHandlerParams)
)

// NewPath returns a new API Path lie /api/something
func NewPath(root, doc string) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).
		Doc(doc).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	return ws
}

// UnsecuredHandle is a handler for an open to the world endpoint
func UnsecuredHandle(handler EndpointFunction) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		handler(
			&EndpointHandlerParams{
				Req:  req,
				Resp: resp,
			},
		)
	}
}
