package plugins

// These are interfaces available from chat server for plugins to use
type (
	// Server hosts chats
	Server interface {
		RemoveUser(nickName string) error
		AddChannel(channelName string)
		GetChannelCount() int
		RemoveUserFromChannel(nickName, channelName string) error
		BroadcastInChannel(channelName, message string) error
		GetChannelUsers(channelName string) (map[string]bool, error)
	}
	// Channel is an interface for a chat channel
	Channel interface {
		AddUser(nickName string)
		GetName() string
		GetUserCount() int
	}
	// User is a connected user
	User interface {
		SetOutgoing(messageType string, message string)
		SetOutgoingf(messageType string, format string, a ...interface{})
		GetChannel() string
		GetNickName() string
		Ignore(nickName string)
		SetChannel(name string)
		HasIgnored(nickName string) bool
		Disconnect() error
	}
)

// These are the different message types that can be sent to user, the purpose of this is to make it easy for consumer to react
const (
	UserOutPutTypeFERunFunction  = "run-function"
	UserOutPutTypeLogInfo        = "log-info"
	UserOutPutTypeLogErr         = "log-error"
	UserOutPutTypeLogWarn        = "log-warning"
	UserOutPutTPM                = "msg"
	UserOutPutTChannel           = "channel"
	UserOutPutTUserTraffic       = "server-notifications"
	UserOutPutTUserTest          = "test"
	UserOutPutTUserInputReq      = "input-req"
	UserOutPutTUserCommandOutput = "command-output"
	UserOutPutTUserServerMessage = "server-message"
)
