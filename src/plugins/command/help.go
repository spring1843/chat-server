package command

import "github.com/spring1843/chat-server/src/shared/errs"

// HelpCommand Shows list of available commands
type HelpCommand struct {
	Command
}

var helpCommand = &HelpCommand{
	Command{
		Name:           `help`,
		Syntax:         `/help`,
		Description:    `Shows the list of all available commands`,
		RequiredParams: []string{`user1`},
	}}

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
