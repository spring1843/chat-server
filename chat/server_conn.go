package chat

import (
	"strconv"

	"github.com/spring1843/chat-server/drivers"
	"github.com/spring1843/chat-server/plugins/logs"
)

// Listen Makes this server start listening to connections, when a user is connected he or she is welcomed
func (s *Server) Listen() {
	go func() {
		for {
			for connection := range s.Connection {
				logs.Infof("connection \t New connection from address=%s", connection.RemoteAddr().String())
				go s.InterviewUser(connection)
			}
		}
	}()
}

// ReceiveConnection is used when there's a new connection
func (s *Server) ReceiveConnection(conn drivers.Connection) {
	s.Connection <- conn
}

// InterviewUser interviews user and allows him to connect after identification
func (s *Server) InterviewUser(conn drivers.Connection) {
	user := NewConnectedUser(s, conn)
	user.SetOutgoing("Welcome to chat server. There are " + strconv.Itoa(s.ConnectedUsersCount()) + " other users on this server. please enter a nickname")

	// wait for user to enter username
	nickName := user.GetIncoming()

	logs.Infof("connection address %q entered user %q", conn.RemoteAddr().String(), nickName)
	for s.IsUserConnected(nickName) {
		user.SetOutgoing("Another user with this nickname is connected to this server, Please enter a different nickname")
		nickName = user.GetIncoming()
	}
	user.SetNickName(nickName)

	s.connectUser(user, conn)
}

func (s *Server) connectUser(user *User, conn drivers.Connection) {
	s.AddUser(user)
	logs.Infof("connection address %s is now nicknamed %q", conn.RemoteAddr().String(), user.GetNickName())
	user.SetOutgoing("Thanks " + user.nickName + ", now please type /join #channel to join a channel or /help to get all commands")
}
