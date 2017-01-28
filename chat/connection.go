package chat

import "net"

// Connection is an interface for a network connection
type Connection interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	RemoteAddr() net.Addr
}
