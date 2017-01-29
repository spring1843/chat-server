package command

import (
	"strconv"

	"github.com/spring1843/chat-server/plugins/errs"
)

// JoinCommand allows user to join a channel
type JoinCommand struct {
	Command
}

var joinCommand = &JoinCommand{
	Command{
		Name:           `join`,
		Syntax:         `/join #channel`,
		Description:    `Join a channel`,
		RequiredParams: []string{`user1`, `channel`},
	}}

// Execute allows a user to join a channel
func (c *JoinCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User1 param is not set")
	}

	channelName := ""
	if params.Channel == nil {
		channelName, err := ParseChannelFromInput(params.RawInput)
		if err != nil {
			return errs.New("Could not parse channel name")
		}
		params.Server.AddChannel(channelName)
	}

	if params.User1.GetChannel() != "" && params.User1.GetChannel() == params.Channel.GetName() {
		return errs.New("You are already in channel #" + params.Channel.GetName())
	}

	channelName = params.Channel.GetName()

	params.Channel.AddUser(params.User1.GetNickName())
	params.User1.SetChannel(channelName)

	// Welcome user to channel
	params.User1.SetOutgoing("There are " + strconv.Itoa(params.Channel.GetUserCount()) + " other users this channel.")

	//Tell others someone's joining
	return params.Server.BroadcastInChannel(channelName, `@`+params.User1.GetNickName()+` just joined channel #`+params.Channel.GetName())
}
