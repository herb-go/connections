package connections

import (
	"net"
)

//RawConnection raw connection interface
type RawConnection interface {
	//Close close connection.
	//Return any error if raised.
	Close() error
	//Send send message to connction.
	//return any error if raised.
	Send([]byte) error
	//MessagesChan connection message chan
	MessagesChan() chan []byte
	//ErrorsChan connection error chan.
	ErrorsChan() chan error
	//RemoteAddr return connection rempte address.
	RemoteAddr() net.Addr
	//C connection close signal chan.
	C() chan int
}

// OutputConnection connection interface which can output messages.
type OutputConnection interface {
	Close() error
	Send([]byte) error
	ID() string
}

//OutputService output service interface
type OutputService interface {
	//Send message to connection by given id.
	Send(id string, msg []byte) error
	//Close close connection by given id.
	Close(id string) error
}

//InputService input service interface
type InputService interface {
	//MessagesChan return connection message chan.
	MessagesChan() chan *Message
	//ErrorsChan return connection error chan.
	ErrorsChan() chan *Error
	//OnCloseEventsChan return connection on close events chan.
	OnCloseEventsChan() chan OutputConnection
	//OnOpenEventsChan return connection on open events chan.
	OnOpenEventsChan() chan OutputConnection
}

//Info connection info struct
type Info struct {
	//ID connection registered id
	ID string
	//Timestamp connection registered timestamp.
	Timestamp int64
}

//Conn connection struct
type Conn struct {
	RawConnection
	Info *Info
}

//New create new connection.
//Return connection created.
func New() *Conn {
	return &Conn{}
}

//ID get connection id
func (c *Conn) ID() string {
	if c.Info != nil {
		return c.Info.ID
	}
	return ""
}

// Message connection message struct
type Message struct {
	//Message raw message
	Message []byte
	//Conn connection which received message
	Conn OutputConnection
}

//Error connection error struct
type Error struct {
	//Error raw error
	Error error
	//Conn connection which raised error.
	Conn OutputConnection
}
