package chat

import (
	"strconv"

	"github.com/spring1843/chat-server/plugins/logs"
)

// Listen Makes this server start listening to connections, when a user is connected he or she is welcomed
func (s *Server) Listen() {
	go func() {
		for {
			for connection := range s.Connection {
				logs.Infof("connection \t New connection from address=%s", connection.RemoteAddr().String())
				go s.WelcomeNewUser(connection)
			}
		}
	}()
}

// ReceiveConnection is used when there's a new connection
func (s *Server) ReceiveConnection(conn Connection) {
	s.Connection <- conn
}

// DisconnectUser disconnects a user from this server
func (s *Server) DisconnectUser(nickName string) error {
	user, err := s.GetUser(nickName)
	if err != nil {
		return err
	}
	return user.Disconnect(s)
}

// WelcomeNewUser shows a welcome message to a new user and makes a new user entity by asking the new user to pick a nickname
func (s *Server) WelcomeNewUser(connection Connection) {
	user := NewConnectedUser(s, connection)
	user.SetOutgoing("Welcome to chat server. There are " + strconv.Itoa(s.ConnectedUsersCount()) + " other users on this server. please enter a nickname")

	nickName := user.GetIncoming()

	for s.IsUserConnected(nickName) {
		user.SetOutgoing("Another user with this nickname is connected to this server, Please enter a different nickname")
		nickName = user.GetIncoming()
	}

	user.SetNickName(nickName)
	s.AddUser(user)
	logs.Infof("connection \t address=%s authenticated=@%s", connection.RemoteAddr().String(), nickName)

	user.SetOutgoing("Thanks " + user.nickName + ", now please type /join #channel to join a channel or /help to get all commands")
}
