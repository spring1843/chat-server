package rest

import (
	"github.com/spring1843/chat-server/libs/go-restful"
	"github.com/spring1843/chat-server/src/shared/rest"
)

// Register the status endpoint
func registerStatusPath(container rest.Container) {
	apiPath := rest.NewPath("/api/status", "Returns the status")

	apiPath.Route(apiPath.GET("").To(Status).
		Writes(statusResp{}))

	container.Add(apiPath)
}

type (
	statusResp struct {
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
