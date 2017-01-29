package command

// Params is all the params that are supported by chat commands, a chat command may use some or all of these params
type Params struct {
	User1    Chatter
	User2    Chatter
	Channel  Chan
	Message  string
	RawInput string
	Server   Server
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
