// Package telnet provides a driver for a chat-server
// When started connections can be made to a tcp port by a telnet like application
package telnet

import (
	"net"

	"github.com/spring1843/chat-server/src/chat"
	"github.com/spring1843/chat-server/src/config"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

// Start starts the telnet server and configures it
func Start(chatServer *chat.Server, config config.Config) error {
	listener, err := net.Listen("tcp", config.TelnetAddress)

	if err != nil {
		logs.ErrIfErrf(err, "error port in use? Error while listening for telnet connections on %s.", config.TelnetAddress)
		return errs.Wrap(err, "Could not open telnet connection")
	}

	go func(chatServer *chat.Server) {
		for {
			connection, err := listener.Accept()
			if err != nil {
				logs.Infof("Error accepting connection %s", err.Error())
			}
			chatServer.ReceiveConnection(&chatConnection{conn: connection})
		}
	}(chatServer)

	return nil
}
