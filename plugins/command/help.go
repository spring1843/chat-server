package command

import "github.com/spring1843/chat-server/plugins/errs"

// Execute shows all available chat commands on this server
func (c *HelpCommand) Execute(params Params) error {
	if params.User1 == nil {
		return errs.New("User param is not set")
	}

	helpMessage := "Here is the list of all available commands\n"
	for _, command := range AllChatCommands {
		helpMessage += command.GetChatCommand().Syntax + "\t" + command.GetChatCommand().Description + ".\n"
	}

	params.User1.SetOutgoing(helpMessage)
	return nil
}
