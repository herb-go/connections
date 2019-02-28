package connections

import (
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

var DefaultIDGenerator = func() (string, error) {
	unid, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return unid.String(), nil
}

func NewGateway() *Gateway {
	return &Gateway{
		IDGenerator:   DefaultIDGenerator,
		messages:      make(chan *Message),
		errors:        make(chan *Error),
		onCloseEvents: make(chan ConnectionOutput),
		onOpenEvents:  make(chan ConnectionOutput),
	}
}

type Gateway struct {
	ID            string
	IDGenerator   func() (string, error)
	Connections   sync.Map
	messages      chan *Message
	errors        chan *Error
	onCloseEvents chan ConnectionOutput
	onOpenEvents  chan ConnectionOutput
}

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
	m.onOpenEvents <- r
	go func() {
		defer func() {
			m.Connections.Delete(r.Info.ID)
		}()
	Listener:
		for {
			select {
			case message := <-conn.Messages():
				m.messages <- &Message{
					Message: message,
					Conn:    r,
				}
			case err := <-conn.Errors():
				m.errors <- &Error{
					Error: err,
					Conn:  r,
				}
			case <-conn.C():
				break Listener
			}
		}
		m.onCloseEvents <- r
	}()
	m.Connections.Store(id, r)
	return r, nil
}
func (m *Gateway) Conn(id string) *Conn {
	val, ok := m.Connections.Load(id)
	if ok == false {
		return nil
	}
	r := val.(*Conn)
	return r
}
func (m *Gateway) Send(id string, msg []byte) error {
	c := m.Conn(id)
	if c == nil {
		return nil
	}
	return c.Send(msg)
}

func (m *Gateway) Close(id string) error {
	c := m.Conn(id)
	if c == nil {
		return nil
	}
	return c.Close()
}

func (m *Gateway) Messages() chan *Message {
	return m.messages
}
func (m *Gateway) Errors() chan *Error {
	return m.errors
}

func (m *Gateway) OnCloseEvents() chan ConnectionOutput {
	return m.onCloseEvents
}
func (m *Gateway) OnOpenEvents() chan ConnectionOutput {
	return m.onOpenEvents
}
