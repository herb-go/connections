package connections

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
)

var currentGenerated = uint32(0)

//DefaultIDGenerator default  generator.
//Return id string and any error if raised.
var DefaultIDGenerator = func() (string, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.BigEndian, time.Now().UnixNano())
	if err != nil {
		return "", err
	}
	c := atomic.AddUint32(&currentGenerated, 1)
	err = binary.Write(buf, binary.BigEndian, c)
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

//NewGateway create new gateway
func NewGateway() *Gateway {
	return &Gateway{
		IDGenerator:   DefaultIDGenerator,
		messages:      make(chan *Message),
		errors:        make(chan *Error),
		onCloseEvents: make(chan OutputConnection),
		onOpenEvents:  make(chan OutputConnection),
		c:             make(chan bool),
		closed:        0,
	}
}

// Gateway connection gateway struct
type Gateway struct {
	//ID gateway id
	ID string
	//IDGenerator connection id generator
	IDGenerator func() (string, error)
	//Connections connections managerd by gateway.
	Connections sync.Map
	//messages connection message chan.
	messages chan *Message
	//errors connection error chan
	errors chan *Error
	//onCloseEvents connection closed event chan
	onCloseEvents chan OutputConnection
	//onOpenEvents connection open event chan
	onOpenEvents chan OutputConnection
	// c close signal chan
	c chan bool
	//closed  if gateway closed
	closed int32
}

// C return close chan.
func (g *Gateway) C() chan bool {
	return g.c
}

//Stop close gate way and close gate close chan.
func (g *Gateway) Stop() {
	if g.isClosed() {
		return
	}
	g.close()
	close(g.C())
	go func() {
		g.Connections.Range(func(key interface{}, value interface{}) bool {
			value.(*Conn).Close()
			return true
		})
	}()
}
func (g *Gateway) isClosed() bool {
	return atomic.LoadInt32(&g.closed) != 0
}
func (g *Gateway) close() {
	atomic.StoreInt32(&g.closed, 1)
}

//Register register raw connection to gateway.
//Return connection and any error if raised.
func (g *Gateway) Register(conn RawConnection) (*Conn, error) {
	id, err := g.IDGenerator()
	if err != nil {
		return nil, err
	}
	if g.ID != "" {
		id = g.ID + "-" + id
	}
	r := &Conn{
		RawConnection: conn,
		Info: &Info{
			ID:        id,
			Timestamp: time.Now().Unix(),
		},
	}

	_, loaded := g.Connections.LoadOrStore(id, r)
	if loaded {
		r.Close()
		return nil, ErrConnIDDuplicated
	}
	go func() {
		if g.isClosed() {
			return
		}
		g.onOpenEvents <- r
	}()
	go func() {
		defer func() {
			g.Connections.Delete(r.Info.ID)
		}()
	Listener:
		for {
			select {
			case message := <-conn.MessagesChan():
				if g.isClosed() {
					return
				}
				g.messages <- &Message{
					Message: message,
					Conn:    r,
				}
			case err := <-conn.ErrorsChan():
				if g.isClosed() {
					return
				}
				g.errors <- &Error{
					Error: err,
					Conn:  r,
				}
			case <-conn.C():
				break Listener
			}
		}
		go func() {
			if g.isClosed() {
				return
			}
			g.onCloseEvents <- r
		}()
	}()
	return r, nil
}

//Conn get connection by id.
//Return nil if connection not registered.
func (g *Gateway) Conn(id string) *Conn {
	val, ok := g.Connections.Load(id)
	if ok == false {
		return nil
	}
	r := val.(*Conn)
	return r
}

//Send send message to connection by given id.
//Return and error if raised.
func (g *Gateway) Send(id string, msg []byte) error {
	c := g.Conn(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}

//Close connection by given id.
//Return any error if raised.
func (g *Gateway) Close(id string) error {
	c := g.Conn(id)
	if c == nil {
		return nil
	}
	g.Connections.Delete(id)
	return c.Close()
}

//ListConn list all connections
func (g *Gateway) ListConn() []*Conn {
	var result = []*Conn{}
	g.Connections.Range(func(key interface{}, value interface{}) bool {
		result = append(result, value.(*Conn))
		return true
	})
	return result
}

//MessagesChan return connection message  chan.
func (g *Gateway) MessagesChan() chan *Message {
	return g.messages
}

//ErrorsChan return connection error chan.
func (g *Gateway) ErrorsChan() chan *Error {
	return g.errors
}

//OnCloseEventsChan return closed event chan.
func (g *Gateway) OnCloseEventsChan() chan OutputConnection {
	return g.onCloseEvents
}

//OnOpenEventsChan return open event chan.
func (g *Gateway) OnOpenEventsChan() chan OutputConnection {
	return g.onOpenEvents
}
