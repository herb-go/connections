package connections

import (
	"net"
	"sync"
)

//NewDummyConnection create new dummy connection.
//Return dummy connection created.
func NewDummyConnection() *DummyConnection {
	return &DummyConnection{
		messages: make(chan []byte, 10),
		Output:   make(chan []byte, 10),
		errors:   make(chan error, 10),
		c:        make(chan int),
	}
}

//DummyConnection dummy connection for testing
type DummyConnection struct {
	Addr     net.Addr
	messages chan []byte
	Output   chan []byte
	errors   chan error
	c        chan int
	Closed   bool
	Lock     sync.Mutex
}

//Close close connection.
//Return any error if raised.
func (c *DummyConnection) Close() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	close(c.c)
	c.Closed = true
	return nil
}

//Send send message to connction.
//return any error if raised.
func (c *DummyConnection) Send(msg []byte) error {
	c.Output <- msg
	return nil
}

//ClientSend Send send message to connction from client.
//return any error if raised.
func (c *DummyConnection) ClientSend(msg []byte) error {
	go func() {
		c.messages <- msg
	}()
	return nil
}

//MessagesChan connection message chan
func (c *DummyConnection) MessagesChan() chan []byte {
	return c.messages
}

//ErrorsChan connection error chan.
func (c *DummyConnection) ErrorsChan() chan error {
	return c.errors
}

//RemoteAddr return connection rempte address.
func (c *DummyConnection) RemoteAddr() net.Addr {
	return c.Addr
}

//C connection close signal chan.
func (c *DummyConnection) C() chan int {
	return c.c
}

//RaiseError raise en error to connection
func (c *DummyConnection) RaiseError(err error) {
	c.errors <- err
}
