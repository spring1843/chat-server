package command

import "github.com/spring1843/chat-server/src/plugins"

// Params is all the params that are supported by chat commands, a chat command may use some or all of these params
type Params struct {
	User1    plugins.Chatter
	User2    plugins.Chatter
	Channel  plugins.Chan
	Message  string
	RawInput string
	Server   plugins.Server
}

// RequiresParam checks to see if a command requires the given parameter
func (c *Command) RequiresParam(param string) bool {
	params := c.GetChatCommand().RequiredParams
	for _, p := range params {
		if p == param {
			return true
		}
	}
	return false
}
