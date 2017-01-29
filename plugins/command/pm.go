package command

import (
	"time"

	"github.com/spring1843/chat-server/plugins/errs"
	"github.com/spring1843/chat-server/plugins/logs"
)

// Execute allows a user to send a private message to another user
func (c *PrivateMessageCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User1 param is not set")
	}

	if params.User2 == nil {
		return errs.New("User2 param is not set")
	}

	if params.User1.GetChannel() != params.User2.GetChannel() {
		return errs.New("Users are not in the same channel")
	}

	if params.User2.HasIgnored(params.User1.GetNickName()) {
		return errs.New("User has ignored the sender")
	}

	logs.Infof("message \t @%s to @%s message=%s", params.User1.GetNickName(), params.User2.GetNickName(), params.Message)

	now := time.Now()
	go params.User2.SetOutgoing(now.Format(time.Kitchen) + ` - *Private from @` + params.User1.GetNickName() + `: ` + params.Message)
	return nil
}
