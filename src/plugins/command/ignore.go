package command

import (
	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/errs"
)

// IgnoreCommand allows a user to ignore another user
type IgnoreCommand struct {
	Command
}

var ignoreCommand = &IgnoreCommand{
	Command{
		Name:           `ignore`,
		Syntax:         `/ignore @nickname`,
		Description:    `Ignore a user, followed by user nickname. An ignored user can not send you private messages`,
		RequiredParams: []string{`user1`, `user2`},
	}}

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
	params.User1.SetOutgoingf(plugins.UserOutPutTUserCommandOutput, "@%s is now ignored.", params.User2.GetNickName())
	return nil
}
