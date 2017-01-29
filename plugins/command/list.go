package command

import (
	"github.com/spring1843/chat-server/plugins/errs"
)

// ListCommand list available users in a channel
type ListCommand struct {
	Command
}

var listCommand = &ListCommand{
	Command{
		Name:           `list`,
		Syntax:         `/list`,
		Description:    `Lists user nicknames in the current channel`,
		RequiredParams: []string{`user1`},
	}}

// Execute lists all the users in a channel to a user
func (c *ListCommand) Execute(params Params) error {
	listMessage := "Here is the list of all users in this channel\n"

	if params.User1 == nil {
		return errs.New("User param is not set")
	}

	channelName := params.User1.GetChannel()
	if channelName == "" {
		return errs.New("User is not in a channel")
	}

	userNickName := params.User1.GetNickName()

	channelUsers, err := params.Server.GetChannelUsers(channelName)
	if err != nil {
		return errs.Wrapf(err, "Error getting channel users")
	}

	for nickName := range channelUsers {
		if nickName == userNickName {
			continue
		}
		listMessage += "@" + nickName + ".\n"
	}
	params.User1.SetOutgoing(listMessage)
	return nil
}
