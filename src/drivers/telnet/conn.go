package telnet

import "net"

type chatConnection struct {
	conn     net.Conn
	nickname string
}

func (c *chatConnection) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

func (c *chatConnection) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *chatConnection) Close() error {
	return c.conn.Close()
}

func (c *chatConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *chatConnection) SetUserNickname(nickName string) {
	c.nickname = nickName
}
