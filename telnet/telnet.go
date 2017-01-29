// Package telnet provides a driver for a chat-server
// When started connections can be made to a tcp port by a telnet like application
package telnet

import (
	"net"
	"strconv"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/config"
	"github.com/spring1843/chat-server/plugins/errs"
	"github.com/spring1843/chat-server/plugins/logs"
)

// Start starts the telnet server and configures it
func Start(chatServer *chat.Server, config config.Config) error {
	listener, err := net.Listen("tcp", config.IP+`:`+strconv.Itoa(config.TelnetPort))

	if err != nil {
		logs.Infof("error \t port in use? Error while listening for telnet connections on %s:%d : %v\n", config.IP, config.TelnetPort, err)
		return errs.Wrap(err, "Could not open telnet connection")
	}

	go func(chatServer *chat.Server) {
		for {
			connection, err := listener.Accept()
			if err != nil {
				logs.Infof("Error accepting connection %s", err.Error())
			}
			chatServer.ReceiveConnection(connection)
		}
	}(chatServer)

	return nil
}
