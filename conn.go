package connections

import (
	"net"
)

type RawConnection interface {
	Close() error
	Send([]byte) error
	Messages() chan []byte
	Errors() chan error
	RemoteAddr() net.Addr
	C() chan int
}

type ConnectionOutput interface {
	Close() error
	Send([]byte) error
	ID() string
}
type ConnectionsOutput interface {
	Send(id string, msg []byte) error
	Close(id string) error
}

type ConnectionsInput interface {
	Messages() chan *Message
	Errors() chan *Error
	OnCloseEvents() chan ConnectionOutput
	OnOpenEvents() chan ConnectionOutput
}

type Info struct {
	ID        string
	Timestamp int64
}

type Conn struct {
	RawConnection
	Info *Info
}

func New() *Conn {
	return &Conn{}
}
func (c *Conn) ID() string {
	if c.Info != nil {
		return c.Info.ID
	}
	return ""
}

type Message struct {
	Message []byte
	Conn    ConnectionOutput
}

type Error struct {
	Error error
	Conn  ConnectionOutput
}
