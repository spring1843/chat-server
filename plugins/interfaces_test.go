package plugins_test

import (
	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/plugins"
)

var _ plugins.Server = chat.NewServer()
var _ plugins.Chan = chat.NewChannel()
var _ plugins.Chatter = chat.NewUser("foo")
