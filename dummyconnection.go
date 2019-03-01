package connections

import "net"

//NewDummyConnection create new dummy connection.
//Return dummy connection created.
func NewDummyConnection() *DummyConnection {
	return &DummyConnection{
		messages: make(chan []byte),
		Output:   make(chan []byte),
		errors:   make(chan error),
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
