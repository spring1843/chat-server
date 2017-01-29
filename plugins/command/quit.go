package command

import "github.com/spring1843/chat-server/plugins/errs"

// QuitCommand allows a user to disconnect from the server
type QuitCommand struct {
	Command
}

var quitCommand = &QuitCommand{
	Command{
		Name:           `quit`,
		Syntax:         `/quit`,
		Description:    `Quit chat server`,
		RequiredParams: []string{`user1`},
	}}

// Execute disconnects a user from server
func (c *QuitCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User1 param is not set")
	}
	if err := params.Server.RemoveUser(params.User1.GetNickName()); err != nil {
		return errs.Wrapf(err, "Could not remove user afeter quit command. . User :%s", params.User1.GetNickName())
	}
	if params.User1.GetChannel() != "" {
		params.Server.RemoveUserFromChannel(params.User1.GetNickName(), params.User1.GetChannel())
	}
	return params.User1.Disconnect()
}
