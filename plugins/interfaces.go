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
	// Chan is an interface for a chat channel
	Chan interface {
		AddUser(nickName string)
		GetName() string
		GetUserCount() int
	}
	// Chatter is a connected user
	Chatter interface {
		SetOutgoing(message string)
		GetChannel() string
		GetNickName() string
		Ignore(nickName string)
		SetChannel(name string)
		HasIgnored(nickName string) bool
		Disconnect() error
	}
)
