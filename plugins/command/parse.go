package command

import (
	"strings"
	"github.com/spring1843/chat-server/plugins/errs"
)


// ParseNickNameFomInput parses a nickname starting with @ from string
func ParseNickNameFomInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "@") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errs.New(`@nickname not found in the input`)
	}
	return result[1:], nil
}

// ParseChannelFromInput parses a channel name starting with # from string
func  ParseChannelFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "#") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errs.New(`#channel not found in the input`)
	}
	return result[1:], nil
}

// ParseCommandFromInput parses a command starting with / from string
func ParseCommandFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "/") >= 0
	})

	result := strings.Join(subStrings, " ")

	if result == `` {
		return ``, errs.New(`/command not found in the input`)
	}
	return result[1:], nil
}

// ParseMessageFromInput parses a message from command, messages do not start with # or / or @
func ParseMessageFromInput(input string) (string, error) {
	subStrings := Filter(input, func(x string) bool {
		return strings.Index(x, "#") != 0 &&
			strings.Index(x, "@") != 0 &&
			strings.Index(x, "/") != 0
	})

	result := strings.Join(subStrings, " ")
	if result == `` {
		return ``, errs.New(`Message not found in the input`)
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
