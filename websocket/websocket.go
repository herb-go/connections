package websocket

import (
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//MsgTypeText text  message type
const MsgTypeText = websocket.TextMessage

//MsgTypeBinary binary  message type
const MsgTypeBinary = websocket.BinaryMessage

// ErrMsgTypeNotMatch error message type not match
var ErrMsgTypeNotMatch = errors.New("websocket message type not match")

// Conn websocket connection
type Conn struct {
	*websocket.Conn
	closed      bool
	messageType int
	closelocker sync.Mutex
	messages    chan []byte
	output      chan []byte
	errors      chan error
	c           chan int
	options     *Options
}

//Send send message to connction.
//return any error if raised.
func (c *Conn) Send(msg []byte) error {
	delay := time.NewTimer(c.options.WriteTimeout)
	c.closelocker.Lock()
	if !c.closed {
		c.closelocker.Unlock()
		select {
		case c.output <- msg:
			delay.Stop()
		case <-delay.C:
		}
	} else {
		c.closelocker.Unlock()
	}
	return nil
}

//C connection close signal chan.
func (c *Conn) C() chan int {
	return c.c
}

//MessagesChan connection message chan
func (c *Conn) MessagesChan() chan []byte {
	return c.messages
}

//ErrorsChan connection error chan.
func (c *Conn) ErrorsChan() chan error {
	return c.errors
}

//Close close connection.
//Return any error if raised.
func (c *Conn) Close() error {
	defer c.closelocker.Unlock()
	c.closelocker.Lock()
	return c.doClose()
}

func (c *Conn) doClose() error {
	if c.closed {
		return nil
	}
	close(c.c)
	c.closed = true
	c.WriteControl(websocket.CloseMessage, nil, c.options.WriteTimeout)
	return c.Conn.Close()
}
func (c *Conn) send(m []byte) error {
	c.closelocker.Lock()
	if c.closed {
		c.closelocker.Unlock()
		return nil
	}
	c.closelocker.Unlock()
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.options.WriteTimeout))
	if err != nil {
		return err
	}
	return c.Conn.WriteMessage(c.messageType, m)
}

//RemoteAddr return connection rempte address.
func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//New create new websocket connection.
//Return connection created.
func New() *Conn {
	return &Conn{
		closed:   true,
		c:        make(chan int),
		messages: make(chan []byte),
		errors:   make(chan error),
		output:   make(chan []byte),
		options:  NewOptions(),
	}
}

const DefaultReadTimeout = time.Minute
const DefaultWriteTimeout = time.Minute

// Upgrader websocket connection upgrader config
var Upgrader = websocket.Upgrader{}

type Options struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MsgType      int
}

func NewOptions() *Options {
	return &Options{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		MsgType:      websocket.TextMessage,
	}
}

//Upgrade Upgrade http requret with given message type to websocket concection.
//Return websocker connection and any error if raised.
func Upgrade(w http.ResponseWriter, r *http.Request, opt *Options) (*Conn, error) {
	wc, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := New()
	if opt != nil {
		c.options = opt
	}
	c.closed = false
	c.Conn = wc
	c.messageType = c.options.MsgType
	go func() {
		for {
			select {
			case m := <-c.output:
				err := c.send(m)
				if err != nil {
					// c.errors <- err
					c.Close()
				}
			case <-c.C():
				// c.closelocker.Lock()
				// close(c.output)
				// close(c.errors)
				// close(c.messages)
				// c.closelocker.Unlock()
				return
			}
		}
	}()

	go func() {

		defer func() {
			recover()
		}()
		for {
			mt, msg, err := c.ReadMessage()
			if err == io.EOF {
				return
			}
			if err != nil {
				c.closelocker.Lock()
				closed := c.closed
				if closed {
					c.closelocker.Unlock()
					return
				}
				if websocket.IsUnexpectedCloseError(err) {
					c.doClose()
					c.closelocker.Unlock()
					return
				}
				c.closelocker.Unlock()
				c.errors <- err
				continue
			}
			if mt != c.messageType {
				c.errors <- ErrMsgTypeNotMatch
				continue
			}
			c.messages <- msg
		}
	}()
	return c, nil
}
