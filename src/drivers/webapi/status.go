package webapi

import "github.com/spring1843/chat-server/src/shared/rest"

// Register the status endpoint
func registerStatusPath(container rest.Container) {
	apiPath := rest.NewPath("/api/status", "Returns the status")
	defer container.Add(apiPath)

	apiPath.Route(apiPath.GET("").To(rest.UnsecuredHandler(getStatus)).
		Operation("getStatus").
		Writes(statusResp{}))

}

type (
	statusResp struct {
		rest.Response
		Data struct {
			Health string `json:"health"`
		}
	}
)

// Status shows the status of the chat server to the users
func getStatus(params *rest.EndpointHandlerParams) {
	resp := new(statusResp)
	resp.Data.Health = `ok`
	resp.DecorateResponse(params.Req)
	params.Resp.WriteEntity(resp)
}
