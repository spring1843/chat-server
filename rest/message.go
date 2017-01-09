package rest

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/emicklei/go-restful"
)

var maxQueryResults = 100

type MessageRequest struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Response
	RecipientUsers []string `json:"recipients"`
}

var MessageNoUsersError = ResponseError{
	Severity:             5,
	HumanFriendlyMessage: `No users are connected to this server`,
	ShortMessage:         `no-connected-users`,
}

var InvalidPatternError = ResponseError{
	Severity:             10,
	HumanFriendlyMessage: `Regex pattern entered is not RE2 compliant`,
	ShortMessage:         `invalid-regex-pattern`,
}

var CouldNotReadLogFileError = ResponseError{
	Severity:             10,
	HumanFriendlyMessage: `Could not read the log file`,
	ShortMessage:         `could-not-read-log`,
}

var TooManyResultsError = ResponseError{
	Severity:             10,
	HumanFriendlyMessage: `Too many results, returning only the first ` + string(maxQueryResults),
	ShortMessage:         `could-not-read-log`,
}

type SearchLogResponse struct {
	Response
	Occurrences []string `json:"occurrences"`
}

func (r RestMessageResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/api").
		Doc("Interact with chat server").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("message").To(r.message).
		Doc("Broadcasts a public announcement to all users connected to the server").
		Reads(MessageRequest{}).
		Writes(MessageResponse{}))

	ws.Route(ws.GET("message").To(r.searchLogForMessages).
		Doc("Searches private and public messages Returns only up to " + string(maxQueryResults) + " messages").
		Param(ws.QueryParameter("pattern", `Optional RE2 Regex pattern to query messages. Examples: '.*' for all logs`).DataType("string")).
		Writes(SearchLogResponse{}))

	container.Add(ws)
}

// RESTful endpoint that broadcasts a message to all users connected to the server
func (r *RestMessageResource) message(request *restful.Request, response *restful.Response) {
	messageResponse := new(MessageResponse)
	messageRequest := new(MessageRequest)
	ParseRequestBody(request, messageRequest)

	if len(r.ChatServer.Users) == 0 {
		messageResponse.AddError(MessageNoUsersError)
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

// RESTful endpoint that searches the log files with a given RE2 regex pattern
func (r *RestMessageResource) searchLogForMessages(request *restful.Request, response *restful.Response) {
	messageResponse := new(SearchLogResponse)
	pattern := request.QueryParameter(`pattern`)

	// if no pattern specified default to all
	if pattern == `` {
		pattern = ".*"
	}

	_, err := regexp.Compile(pattern)
	if err != nil {
		InvalidPatternError.HumanFriendlyMessage = "Invalid REGEX pattern, pattern is not RE2 compliant:" + pattern
		messageResponse.AddError(InvalidPatternError)
		response.WriteEntity(messageResponse)
		return
	}

	file, err := os.Open(LogFilePath)
	defer file.Close()

	if err != nil {
		messageResponse.AddError(CouldNotReadLogFileError)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	occurrences := make([]string, 0)
	regex := regexp.MustCompile(pattern)
	for scanner.Scan() {

		if len(occurrences) == 100 {
			messageResponse.AddError(TooManyResultsError)
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
