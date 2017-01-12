package rest

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/emicklei/go-restful"
)

// Register the message REST endpoints
func (r messageEndpoint) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/api").
		Doc("Interact with chat server").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("message").To(r.broadCastMessage).
		Doc("Broadcasts a public announcement to all users connected to the server").
		Reads(messageReq{}).
		Writes(messageResp{}))

	ws.Route(ws.GET("message").To(r.searchLogForMessages).
		Doc("Searches private and public messages Returns only up to " + string(maxQueryResults) + " messages").
		Param(ws.QueryParameter("pattern", `Optional RE2 Regex pattern to query messages. Examples: '.*' for all logs`).DataType("string")).
		Writes(searchLogResp{}))

	container.Add(ws)
}

type (
	messageReq struct {
		Message string `json:"message"`
	}
	messageResp struct {
		Response
		RecipientUsers []string `json:"recipients"`
	}
	searchLogResp struct {
		Response
		Occurrences []string `json:"occurrences"`
	}
)

var (
	maxQueryResults   = 100
	errMessageNoUsers = ResponseError{
		Severity:             5,
		HumanFriendlyMessage: `No users are connected to this server`,
		ShortMessage:         `no-connected-users`,
	}
	errInvalidPattern = ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Regex pattern entered is not RE2 compliant`,
		ShortMessage:         `invalid-regex-pattern`,
	}
	errCouldNotReadLogFile = ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Could not read the log file`,
		ShortMessage:         `could-not-read-log`,
	}
	errTooManyResults = ResponseError{
		Severity:             10,
		HumanFriendlyMessage: `Too many results, returning only the first ` + string(maxQueryResults),
		ShortMessage:         `could-not-read-log`,
	}
)

func (r *messageEndpoint) broadCastMessage(request *restful.Request, response *restful.Response) {
	messageResponse := new(messageResp)
	messageRequest := new(messageReq)
	ParseRequestBody(request, messageRequest)

	if len(r.ChatServer.Users) == 0 {
		messageResponse.AddError(errMessageNoUsers)
	}

	r.ChatServer.LogPrintf("message \t RESTful public annoucnement=%s", messageRequest.Message)

	var recipientUsers []string
	for _, user := range r.ChatServer.Users {
		user.Outgoing <- "Public Server Announcement: " + messageRequest.Message
		recipientUsers = append(recipientUsers, "@"+user.NickName)
	}

	messageResponse.RecipientUsers = recipientUsers
	messageResponse.DecorateResponse(request)
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
