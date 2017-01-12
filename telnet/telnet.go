// Package telnet provides a driver for a chat-server
// When started connections can be made to a tcp port by a telnet like application
package telnet

import (
	"errors"
	"net"
	"strconv"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/integration"
)

// Start starts the telnet server and configures it
func Start(chatServer *chat.Server, config integration.Config) error {
	listener, err := net.Listen("tcp", config.IP+`:`+strconv.Itoa(config.TelnetPort))

	if err != nil {
		chatServer.LogPrintf("error \t port in use? Error while listening for telnet connections on %s:%d : %v\n", config.IP, config.TelnetPort, err)
		return errors.New("Could not open telnet connection")
	}

	go func(chatServer *chat.Server) {
		for {
			connection, err := listener.Accept()
			if err != nil {
				chatServer.Logger.Printf("Error accepting connection %s", err.Error())
			}
			chatServer.Connection <- connection
		}
	}(chatServer)

	return nil
}
