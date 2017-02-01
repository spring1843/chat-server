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
		SetOutgoing(message string)
		SetOutgoingf(format string, a ...interface{})
		GetChannel() string
		GetNickName() string
		Ignore(nickName string)
		SetChannel(name string)
		HasIgnored(nickName string) bool
		Disconnect() error
	}
)
