package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"errors"

	"github.com/spring1843/chat-server/chat"
	"github.com/spring1843/chat-server/integration"
	"github.com/spring1843/chat-server/rest"
	"github.com/spring1843/chat-server/telnet"
	"github.com/spring1843/chat-server/websocket"
)

func main() {
	if validateCommandArguments(os.Args) == false {
		panic(errors.New("Invalid command arguments"))
	}

	config := getConfig(os.Args[2])
	chatServer := configureNewChatServer(config)
	startTelnet(config, chatServer)
	fmt.Printf("Info - Telnet Server Started\n")
	startWebSocket(config, chatServer)
	fmt.Printf("Info - Websocket Server Started\n")
	startRest(config, chatServer)
	fmt.Printf("Info - Web Server Started\n")

	fmt.Printf("Waiting for TCP connections...\n")
	neverEnd()
}

// configureNewChatServer configures a new chat server
func configureNewChatServer(config integration.Config) *chat.Server {
	chatServer := chat.NewServer()
	chatServer.Listen()
	setLogFile(config.LogFile, chatServer)
	return chatServer
}

// startTelnet starts the telnet server
func startTelnet(config integration.Config, chatServer *chat.Server) {
	err := telnet.Start(chatServer, config)
	if err != nil {
		fmt.Printf("Could not start telnet server please check the logs for more info\n")
		panic(err)
	}
}

// startRest starts the rest server
func startRest(config integration.Config, chatServer *chat.Server) {
	rest.Start(chatServer, config)
}

// startWebSocket starts the Websocket server
func startWebSocket(config integration.Config, chatServer *chat.Server) {
	err := websocket.Start(chatServer, config)
	if err != nil {
		fmt.Printf("Could not start websocket server please check the logs for more info\n")
		panic(err)
	}
}

// validateCommandArguments validates command line arguments
func validateCommandArguments(args []string) bool {
	if len(args) < 3 || args[1] != `-c` || args[2] == `` {
		fmt.Printf("Error - No config file specified. Usage %s -c config.json\n", os.Args[0])
		return false
	}
	return true
}

// setLogFile creates a log file to be given to a server
func setLogFile(logFile string, chatServer *chat.Server) {
	if logFile == `` {
		return
	}

	logger, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Error - opening log file %s : %v\n", logFile, err)
		panic(err)
	}
	chatServer.SetLogFile(logger)
	fmt.Printf("Info - Log files written to %s\n", logFile)
}

// getConfig parses configurations from a json string
func getConfig(configFile string) integration.Config {
	fileContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		panic(err)
	}

	config := new(integration.Config)
	err = json.Unmarshal([]byte(fileContents), &config)
	if err != nil {
		fmt.Printf("Error parsing JSON config file: %v\n", err)
		panic(err)
	}

	return *config
}

// neverEnd never end execution, avoids termination while the server is concurrently running
func neverEnd() {
	for {
		time.Sleep(1000 * time.Second)
	}
}
