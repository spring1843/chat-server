package webapi

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/spring1843/chat-server/libs/go-restful"
	"github.com/spring1843/chat-server/src/shared/logs"
	"github.com/spring1843/chat-server/src/shared/rest"
)

// Register the message REST endpoints
func (r messageEndpoint) Register(container *restful.Container) {
	apiPath := rest.NewPath("/api/message", "Interact with chat server")

	apiPath.Route(apiPath.POST("").To(r.broadCastMessage).
		Doc("Broadcasts a public announcement to all users connected to the server").
		Reads(messageReq{}).
		Writes(messageResp{}))

	apiPath.Route(apiPath.GET("").To(r.searchLogForMessages).
		Doc("Searches private and public messages Returns only up to " + string(maxQueryResults) + " messages").
		Param(apiPath.QueryParameter("pattern", `Optional RE2 Regex pattern to query messages. Examples: '.*' for all logs`).DataType("string")).
		Writes(searchLogResp{}))

	container.Add(apiPath)
}

type (
	messageReq struct {
		Message string `json:"message"`
	}
	messageResp struct {
		rest.Response
		Success bool `json:"success"`
	}
	searchLogResp struct {
		rest.Response
		Occurrences []string `json:"occurrences"`
	}
)

var (
	maxQueryResults   = 100
	errMessageNoUsers = rest.ResponseError{
		Severity:             5,
		HumanFriendlyMessage: `No users are connected to this server`,
		ShortMessage:         `no-connected-users`,
	}
	errInvalidPattern = rest.ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Regex pattern entered is not RE2 compliant`,
		ShortMessage:         `invalid-regex-pattern`,
	}
	errCouldNotReadLogFile = rest.ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Could not read the log file`,
		ShortMessage:         `could-not-read-log`,
	}
	errTooManyResults = rest.ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Too many results, returning only the first ` + string(maxQueryResults),
		ShortMessage:         `could-not-read-log`,
	}
)

func (r *messageEndpoint) broadCastMessage(request *restful.Request, response *restful.Response) {
	messageResponse := new(messageResp)
	messageRequest := new(messageReq)
	rest.ParseRequestBody(request, messageRequest)

	if r.ChatServer.ConnectedUsersCount() == 0 {
		messageResponse.AddError(errMessageNoUsers)
	}

	logs.Infof("message \t RESTful public announcement=%s", messageRequest.Message)

	r.ChatServer.Broadcast("Public Server Announcement: " + messageRequest.Message)

	messageResponse.DecorateResponse(request)
	messageResponse.Success = true
	response.WriteEntity(messageResponse)
}

func (r *messageEndpoint) searchLogForMessages(request *restful.Request, response *restful.Response) {
	messageResponse := new(searchLogResp)
	pattern := request.QueryParameter(`pattern`)

	// if no pattern specified default to all
	if pattern == `` {
		pattern = ".*"
	}

	_, err := regexp.Compile(pattern)
	if err != nil {
		errInvalidPattern.HumanFriendlyMessage = "Invalid REGEX pattern, pattern is not RE2 compliant:" + pattern
		messageResponse.AddError(errInvalidPattern)
		response.WriteEntity(messageResponse)
		return
	}

	file, err := os.Open(LogFilePath)
	defer file.Close()

	if err != nil {
		messageResponse.AddError(errCouldNotReadLogFile)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	var occurrences []string
	regex := regexp.MustCompile(pattern)
	for scanner.Scan() {
		if len(occurrences) == 100 {
			messageResponse.AddError(errTooManyResults)
			break
		}

		token := scanner.Text()
		if strings.Contains(token, "message \t") == false {
			continue
		}

		if entry := regex.FindString(token); entry != `` {
			occurrences = append(occurrences, entry)
		}
	}

	messageResponse.Occurrences = occurrences
	messageResponse.DecorateResponse(request)
	response.WriteEntity(messageResponse)
}
