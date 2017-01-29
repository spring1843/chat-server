package command

type (
	// Executable Is a what a command is
	Executable interface {
		Execute(params Params) error
		GetChatCommand() Command
		ParseNickNameFomInput(input string) (string, error)
		ParseChannelFromInput(input string) (string, error)
		ParseMessageFromInput(input string) (string, error)
		ParseCommandFromInput(input string) (string, error)
	}
	// Server hosts chats
	Server interface {
		RemoveUser(nickName string) error

		AddChannel(channelName string)
		GetChannelCount() int
		RemoveUserFromChannel(nickName, channelName string) error
		BroadcastInChannel(channelName, message string) error
		DisconnectUser(nickName string) error
	}
	// Chan is an interface for a chat channel
	Chan interface {
		AddUser(nickName string)
		GetName() string
		GetUserCount() int
		GetUsers() map[string]bool
	}
	// Chatter is a connected user
	Chatter interface {
		SetOutgoing(message string)
		GetChannel() string
		GetNickName() string
		Ignore(nickName string)
		SetChannel(name string)
		HasIgnored(nickName string) bool
	}
)
