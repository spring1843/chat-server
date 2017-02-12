package rest

import "github.com/spring1843/chat-server/libs/go-restful"

// Container interfaces *restful.Container
type Container interface {
	Add(service *restful.WebService) *restful.Container
}

// NewPath returns a new API Path lie /api/something
func NewPath(root, doc string) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).
		Doc(doc).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	return ws
}
