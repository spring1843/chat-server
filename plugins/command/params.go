package command

import (
	"errors"
	"strings"
)

// Params is all the params that are supported by chat commands, a chat command may use some or all of these params
type Params struct {
	User1    Chatter
	User2    Chatter
	Channel  Chan
	Message  string
	RawInput string
	Server   Server
}

// ParseNickNameFomInput parses a nickname starting with @ from string
func (c *Command) ParseNickNameFomInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "@") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`@nickname not found in the input`)
	}
	return result[1:], nil
}

// ParseChannelFromInput parses a channel name starting with # from string
func (c *Command) ParseChannelFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "#") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`#channel not found in the input`)
	}
	return result[1:], nil
}

// ParseCommandFromInput parses a command starting with / from string
func (c *Command) ParseCommandFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "/") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`/command not found in the input`)
	}
	return result[1:], nil
}

// ParseMessageFromInput parses a message from command, messages do not start with # or / or @
func (c *Command) ParseMessageFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "#") != 0 &&
			strings.Index(x, "@") != 0 &&
			strings.Index(x, "/") != 0
	})

	result := strings.Join(subStrings, " ")
	if result == `` {
		return ``, errors.New(`Message not found in the input`)
	}
	return result, nil
}

// Filter filters input based on the given requirements defined by f function
func Filter(input string, function func(string) bool) []string {
	var vsf []string
	for _, v := range strings.Split(input, " ") {
		if function(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
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
