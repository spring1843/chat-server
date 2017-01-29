package command

import "github.com/spring1843/chat-server/plugins/errs"

// Execute lists all the users in a channel to a user
func (c *ListCommand) Execute(params Params) error {
	listMessage := "Here is the list of all users in this channel\n"

	if params.User1 == nil {
		return errs.New("User param is not set")
	}

	if params.User1.GetChannel() == "" {
		return errs.New("User is not in a channel")
	}

	for nickName := range params.Channel.GetUsers() {
		if nickName == params.User1.GetNickName() {
			continue
		}
		listMessage += "@" + nickName + ".\n"
	}

	params.User1.SetOutgoing(listMessage)

	return nil
}
