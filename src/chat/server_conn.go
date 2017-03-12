package chat

import (
	"github.com/spring1843/chat-server/src/drivers"
	"github.com/spring1843/chat-server/src/plugins"
	"github.com/spring1843/chat-server/src/shared/errs"
	"github.com/spring1843/chat-server/src/shared/logs"
)

// Listen Makes this server start listening to connections, when a user is connected he or she is welcomed
func (s *Server) Listen() {
	go func() {
		for {
			for connection := range s.conn {
				logs.Infof("New connection from address=%s", connection.RemoteAddr().String())
				go s.InterviewUser(connection)
			}
		}
	}()
}

// ReceiveConnection is used when there's a new connection
func (s *Server) ReceiveConnection(conn drivers.Connection) {
	s.conn <- conn
}

// InterviewUser interviews user and allows him to connect after identification
func (s *Server) InterviewUser(conn drivers.Connection) {
	user := NewConnectedUser(conn)
	user.Listen(s)
	if err := s.authenticateUser(user, conn); err != nil {
		logs.ErrIfErrf(err, "Failed authenticating user.")
		logs.ErrIfErrf(conn.Close(), "Failed closing connection after failed authentication")
		return // Don't attach to server
	}
	s.connectUser(user, conn)
}

func (s *Server) authenticateUser(user *User, conn drivers.Connection) error {
	user.SetOutgoingf(plugins.UserOutPutTUserInputReq, "Welcome to chat server. There are %d other users on this server. please enter a nickname", s.ConnectedUsersCount())

	// wait for user to enter username
	nickName, err := user.GetIncoming()
	if err != nil {
		return errs.Wrapf(err, "Failed getting nickname from user connection")
	}

	logs.Infof("connection address %q entered user %q", conn.RemoteAddr().String(), nickName)
	for s.IsUserConnected(nickName) {
		user.SetOutgoingf(plugins.UserOutPutTypeLogErr, "Another user with nickname %q is connected to this server, Please enter a different nickname", nickName)
		nickName, err = user.GetIncoming()
		if err != nil {
			return errs.Wrapf(err, "Failed getting nickname from user connection after picking duplicated nickname %q", nickName)
		}
	}
	user.SetNickName(nickName)
	return nil
}

func (s *Server) connectUser(user *User, conn drivers.Connection) {
	s.AddUser(user)
	nickName := user.GetNickName()
	logs.Infof("connection address %s is now nicknamed %q", conn.RemoteAddr().String(), nickName)
	conn.SetUserNickname(nickName)
	user.SetOutgoingf(plugins.UserOutPutTUserServerMessage, "Welcome %s! Enter /join #channel or /help to get all commands", user.GetNickName())
}
