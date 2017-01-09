package chat

import (
	"errors"
	"strings"
)

// All the params that are supported by chat commands, a chat command may use some or all of these params
type CommandParams struct {
	user1    *User
	user2    *User
	channel  *Channel
	message  string
	rawInput string
}

// Parses a nickname starting with @ from string
func (c *ChatCommand) ParseNickNameFomInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "@") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`@nickname not found in the input`)
	}
	return result[1:], nil
}

// Parses a channel name starting with # from string
func (c *ChatCommand) ParseChannelFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "#") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`#channel not found in the input`)
	}
	return result[1:], nil
}

// Parses a command starting with / from string
func (c *ChatCommand) ParseCommandFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "/") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errors.New(`/command not found in the input`)
	}
	return result[1:], nil
}

// Parses a message from command, messages do not start with # or / or @
func (c *ChatCommand) ParseMessageFromInput(input string) (string, error) {
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

// Filter input based on the given requirements defined by f function
func Filter(input string, function func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range strings.Split(input, " ") {
		if function(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Checks to see if a command requires the given parameter
func (c *ChatCommand) DoesCommandRequireParam(param string) bool {
	params := c.getChatCommand().RequiredParams
	for _, p := range params {
		if p == param {
			return true
		}
	}
	return false
}
