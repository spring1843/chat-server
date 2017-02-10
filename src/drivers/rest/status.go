package rest

import (
	"github.com/spring1843/chat-server/libs/go-restful"
)

// Register the status endpoint
func (g statusEndpoint) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/api/status").
		Doc("Returns the status").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(Status).
		Writes(statusResp{}))

	container.Add(ws)
}

type (
	statusEndpoint struct{}
	statusResp     struct {
		Response
		Data struct {
			Health string `json:"health"`
		}
	}
)

// Status shows the status of the chat server to the users
func Status(request *restful.Request, response *restful.Response) {
	resp := new(statusResp)
	resp.Data.Health = `ok`
	resp.DecorateResponse(request)
	response.WriteEntity(resp)
}
