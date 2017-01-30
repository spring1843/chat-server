package plugins_test

import (
	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/plugins"
)

var _ plugins.Server = chat.NewServer()
var _ plugins.Chan = chat.NewChannel()
var _ plugins.Chatter = chat.NewUser("foo")
