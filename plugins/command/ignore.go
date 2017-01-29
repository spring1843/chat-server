package command

import "github.com/spring1843/chat-server/plugins/errs"

// Execute allows a user to ignore another user so to suppress all incoming messages from that user
func (c *IgnoreCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User1 param is not set")
	}

	if params.User2 == nil {
		return errs.New("User2 param is not set")
	}

	if params.User2.GetNickName() == params.User1.GetNickName() {
		return errs.New("We don't let one ignore themselves")
	}

	params.User1.Ignore(params.User2.GetNickName())
	params.User1.SetOutgoing(params.User2.GetNickName() + " is now ignored.")
	return nil
}
