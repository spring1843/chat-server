package command

import "github.com/spring1843/chat-server/plugins/errs"

// Execute disconnects a user from server
func (c *QuitCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User1 param is not set")
	}

	if err := params.Server.RemoveUser(params.User1.GetNickName()); err != nil {
		return errs.New("Could not remove user afeter quit command")
	}
	if params.User1.GetChannel() != "" {
		params.Server.RemoveUserFromChannel(params.User1.GetNickName(), params.User1.GetChannel())
	}

	if err := params.Server.DisconnectUser(params.User1.GetNickName()); err != nil {
		return errs.New("Could not disconnect user")
	}
	return nil
}
