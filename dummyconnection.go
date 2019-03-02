package connections

import (
	"net"
	"time"
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
}

//Close close connection.
//Return any error if raised.
func (c *DummyConnection) Close() error {
	close(c.c)
	close(c.errors)
	close(c.messages)
	close(c.Output)
	return nil
}

//Send send message to connction.
//return any error if raised.
func (c *DummyConnection) Send(msg []byte) error {
	c.Output <- msg
	return nil
}

//ReadOutput read  connection output with timeout.
func (c *DummyConnection) ReadOutput() ([]byte, bool) {
	select {
	case m, closed := <-c.Output:
		return m, closed
	case <-time.NewTimer(time.Millisecond).C:
		return nil, false
	}
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

//ReadC read connection close chan
func (c *DummyConnection) ReadC() (int, bool) {
	select {
	case v, closed := <-c.C():
		return v, closed
	case <-time.NewTimer(time.Millisecond).C:
		return -1, false
	}
}
