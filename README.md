[![Travis-CI](https://travis-ci.org/spring1843/chat-server.svg?branch=master)](https://travis-ci.org/spring1843/chat-server/) [![Report card](https://goreportcard.com/badge/github.com/spring1843/chat-server)](https://goreportcard.com/report/github.com/spring1843/chat-server)


Chat Server is basic IRC like server written in Go with drivers for Telnet, WebSocket, and HTTP REST.

--
## Telnet and WebSocket Servers
Users can connect to chat server with a simple TCP client like telnet or WebSockets. After connecting to the server users can join or create channels and publicly chat with everyone in the channel. Users can use the /msg command to send private messages to each other. Chat Commands are supported for both Telnet and WebSocket clients.

## Supported Chat Commands
    /help                       Shows the list of all available commands.
    /list                       Lists user nicknames in the current channel.
    /ignore @nickname	        Ignore a user so that he or she can't send you private messages.
    /join #channel              Join a channel.
    /msg @nickname message	    Send a private message to a user.
    /quit	                    Quit chat server.

## HTTP Rest Server
Administrators can use the HTTP RESTful endpoints to broadcast a public announcement to everyone connected to the server. There's also an end point to query the messages stored in log files. To see the swagger documentation for RESTful API's please browse ```/docs``` on the HTTP port (4001 by default).  

## How to Run Locally
- Download and install dependencies by running ```go get github.com/emicklei/go-restful``` and ```go get github.com/gorilla/websocket```
- Download and install by running  ```go get github.com/spring1843/chat-server/src/``` and then CD into the directory
- Run all tests ```make test```
- Modify or review default settings by editing ```config.json```
- To start serving run ```make serve``` or ```go run main.go -c config.json``` or ```chat-server -c config.json```
- To check-in new code run ```make checkin``` to automatically format, lint, vet, race check, run all tests and review each change

## How to Connect
To connect with telnet when the server is running locally try ```telnet localhost 4000```
To connect with a WebSocket open a browser with WebSocket support (such as Google Chrome) and browse to ```http://localhost:4004/client```
These ports are defined in config.json

## Logs
Besides start up messages, everything else is logged in a log file. Almost all interactions including the followings are added to the log file specified at run time along with a timestamp:
- Public Messages
- Private Messages
- Chat commands executed
- New and closing connections
- Errors

## Docker
To build
```
docker build -t chatserver .
```

To run
```
docker run -d -p 4000:4000 -p 4001:4001 -p 4004:4004 -t chatserver
```

Note that if you are running docker in a VM or boot2docker, the host needs to forward these ports. 

## Design

A config structure is initially parsed from the ```.json``` file specified by the ```-c``` flag and passed around to every package.  

The ```chat``` package is in charge of creating the chat experience. Important entities in this package are Server, User, Channel, ChatConnection, and ChatCommand.

Chat server can serve anything that implements the ChatConnection interface providing basic ```io.Writer``` and ```io.Reader``` methods and a few network specific concepts like ```RemoteAddr``` and ```Close```. Examples of such drivers in use are Telnet and WebSockets. 

The Server also exposes certain internal entities and functionality to be used by external entities that do not necessarily have a persistent connection for example the main file is able to attach a logger to the server, and the REST driver can query the server.
 
One instance of Chat server is created at the beginning of execution and passed around to different drivers when they are started. When a new connection is made the Server welcomes the new user and tries to identify the user. Upon successful identification the user and the open connection associated with him are attached to the Server until the connection is closed. 

After identification users can use ```ChatCommand```s are commands that a user can execute on the server. Each command has a syntax and certain parameters. The ```chat``` package is able to identify, parse their parameters from the input and execute them. 

Users can ```/join``` a channel. Users can join only one Channel at any given time and participate in the conversations. They can also send private messages to each other with ```/msg``` command.

Channels, Users and ChatCommands are represented with ```#,@,/``` prefixes respectively.

## Credits
* [go-restful](https://github.com/emicklei/go-restful) for configuring, standardizing, documenting and running HTTP REST Endpoints
* [Swagger](http://swagger.io/) for auto generating HTTP REST API's and providng an easy to use front end client ```http://localhost:4001/docs```
* [Gorilla WebSocket](https://github.com/gorilla/websocket) for handling WebSocket connections and example client ```http://localhost:4004/client```

### Further Improvements
- Strive for 100% test coverage
- Deploy to Heroku
- Handle client tcp disconnect may be ping and pong?
- Use flags package for parsing command line options
