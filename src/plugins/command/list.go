package command

import (
	"strings"

	"github.com/spring1843/chat-server/src/shared/errs"
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
	if params.User1 == nil {
		return errs.New("User param is not set")
	}

	channelName := params.User1.GetChannel()
	if channelName == "" {
		return errs.New("User is not in a channel")
	}

	channelUsers, err := params.Server.GetChannelUsers(channelName)
	if err != nil {
		return errs.Wrapf(err, "Error getting channel users")
	}

	users := make([]string, len(channelUsers), len(channelUsers))
	for nickName := range channelUsers {
		users = append(users, "@"+nickName)
	}
	params.User1.SetOutgoingf("User(s) in #%s: %s", channelName, strings.Join(users, ","))
	return nil
}
