package rest

import "github.com/spring1843/chat-server/libs/go-restful"

func NewPath(root, doc string) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(root).
		Doc(doc).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	return ws
}
