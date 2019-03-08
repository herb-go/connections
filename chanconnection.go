package connections

import (
	"net"
	"sync"
)

//NewChanConnection create new chan connection.
//Return chan connection created.
func NewChanConnection() *ChanConnection {
	return &ChanConnection{
		messages:       make(chan []byte, 10),
		clientmessages: make(chan []byte, 10),
		errors:         make(chan error, 10),
		c:              make(chan int),
	}
}

//ChanConnection chan connection
type ChanConnection struct {
	Addr           net.Addr
	messages       chan []byte
	clientmessages chan []byte
	errors         chan error
	c              chan int
	Closed         bool
	Lock           sync.Mutex
}

// Client create chan connection client.
func (c *ChanConnection) Client() *ChanConnectionClient {
	return &ChanConnectionClient{
		Conn: c,
	}
}

//Close close connection.
//Return any error if raised.
func (c *ChanConnection) Close() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	close(c.c)
	c.Closed = true
	return nil
}

//ClientClose close connection.
//Return any error if raised.
func (c *ChanConnection) ClientClose() error {
	return c.Close()
}

//Send send message to connction.
//return any error if raised.
func (c *ChanConnection) Send(msg []byte) error {
	c.ClientMessagesChan() <- msg
	return nil
}

//ClientSend Send  message to connction from client.
//return any error if raised.
func (c *ChanConnection) ClientSend(msg []byte) error {
	go func() {
		c.messages <- msg
	}()
	return nil
}

//ClientMessagesChan connection client message chan
func (c *ChanConnection) ClientMessagesChan() chan []byte {
	return c.clientmessages
}

//MessagesChan connection message chan
func (c *ChanConnection) MessagesChan() chan []byte {
	return c.messages
}

//ErrorsChan connection error chan.
func (c *ChanConnection) ErrorsChan() chan error {
	return c.errors
}

//RemoteAddr return connection rempte address.
func (c *ChanConnection) RemoteAddr() net.Addr {
	return c.Addr
}

//C connection close signal chan.
func (c *ChanConnection) C() chan int {
	return c.c
}

//RaiseError raise an error to connection
func (c *ChanConnection) RaiseError(err error) {
	c.errors <- err
}

//ChanConnectionClient client for chan connection.
type ChanConnectionClient struct {
	Conn *ChanConnection
}

//MessagesChan connection client message chan
func (c *ChanConnectionClient) MessagesChan() chan []byte {
	return c.Conn.ClientMessagesChan()
}

//Send send message to connction from client.
//return any error if raised.
func (c *ChanConnectionClient) Send(msg []byte) error {
	return c.Conn.ClientSend(msg)
}

//Close close connection.
//Return any error if raised.
func (c *ChanConnectionClient) Close() error {
	return c.Conn.ClientClose()
}
