package drivers

import "net"

// Connection is an interface for a network connection of any kind
// A chat server driver must be able to support these network operations
type Connection interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
	RemoteAddr() net.Addr
}
