package rest

import (
	"github.com/emicklei/go-restful"
)

type StatusResource struct{}

type StatusResponse struct {
	Response
	Data struct {
		Health string `json:"health"`
	}
}

func (g StatusResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/status").
		Doc("Returns the status").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(Status).
		Writes(StatusResponse{}))

	container.Add(ws)
}

func Status(request *restful.Request, response *restful.Response) {
	resp := new(StatusResponse)
	resp.Data.Health = `ok`
	resp.DecorateResponse(request)
	response.WriteEntity(resp)
}
