package connections

import (
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

//DefaultIDGenerator default  generator.
//Return uuid string and any error if raised.
var DefaultIDGenerator = func() (string, error) {
	unid, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return unid.String(), nil
}

//NewGateway create new gateway
func NewGateway() *Gateway {
	return &Gateway{
		IDGenerator:   DefaultIDGenerator,
		messages:      make(chan *Message),
		errors:        make(chan *Error),
		onCloseEvents: make(chan OutputConnection),
		onOpenEvents:  make(chan OutputConnection),
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
}

//Register register raw connection to gateway.
//Return connection and any error if raised.
func (m *Gateway) Register(conn RawConnection) (*Conn, error) {
	id, err := m.IDGenerator()
	if err != nil {
		return nil, err
	}
	if m.ID != "" {
		id = m.ID + "-" + id
	}
	r := &Conn{
		RawConnection: conn,
		Info: &Info{
			ID:        id,
			Timestamp: time.Now().Unix(),
		},
	}
	go func() {
		m.onOpenEvents <- r
	}()
	go func() {
		defer func() {
			m.Connections.Delete(r.Info.ID)
		}()
	Listener:
		for {
			select {
			case message := <-conn.MessagesChan():
				m.messages <- &Message{
					Message: message,
					Conn:    r,
				}
			case err := <-conn.ErrorsChan():
				m.errors <- &Error{
					Error: err,
					Conn:  r,
				}
			case <-conn.C():
				break Listener
			}
		}
		go func() {
			m.onCloseEvents <- r
		}()
	}()
	m.Connections.Store(id, r)
	return r, nil
}

//Conn get connection by id.
//Return nil if connection not registered.
func (m *Gateway) Conn(id string) *Conn {
	val, ok := m.Connections.Load(id)
	if ok == false {
		return nil
	}
	r := val.(*Conn)
	return r
}

//Send send message to connection by given id.
//Return and error if raised.
func (m *Gateway) Send(id string, msg []byte) error {
	c := m.Conn(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}

//Close connection by given id.
//Return any error if raised.
func (m *Gateway) Close(id string) error {
	c := m.Conn(id)
	if c == nil {
		return nil
	}
	return c.Close()
}

//MessagesChan return connection message  chan.
func (m *Gateway) MessagesChan() chan *Message {
	return m.messages
}

//ErrorsChan return connection error chan.
func (m *Gateway) ErrorsChan() chan *Error {
	return m.errors
}

//OnCloseEventsChan return closed event chan.
func (m *Gateway) OnCloseEventsChan() chan OutputConnection {
	return m.onCloseEvents
}

//OnOpenEventsChan return open event chan.
func (m *Gateway) OnOpenEventsChan() chan OutputConnection {
	return m.onOpenEvents
}
