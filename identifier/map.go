package identifier

import (
	"sync"

	"github.com/herb-go/connections"
)

// GenerateDefaultMapOnLogout generate on logout callback  by  given mapidentifier.
var GenerateDefaultMapOnLogout = func(m *Map) func(id string, conn connections.OutputConnection) error {
	return func(id string, conn connections.OutputConnection) error {
		conn.Close()
		return nil
	}
}

//Map identifier  useing sync.Map
type Map struct {
	Identities sync.Map
	lock       sync.Mutex
	onLogout   func(id string, conn connections.OutputConnection) error
}

func (m *Map) conn(id string) (connections.OutputConnection, bool) {
	data, ok := m.Identities.Load(id)
	if ok == false {
		return nil, false
	}
	conn, ok := data.(*connections.Conn)
	return conn, ok
}

// Login given connection as user by given id.
//Return any error if raised.
func (m *Map) Login(id string, c connections.OutputConnection) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	conn, ok := m.conn(id)
	if ok {
		err := m.onLogout(id, conn)
		if err != nil {
			return err
		}
	}
	m.Identities.Delete(id)
	m.Identities.Store(id, c)
	return nil
}

//Logout logout user by given id if current user is given connection.
//Always logout user if given connection is nil.
//Return any error if raised.
func (m *Map) Logout(id string, c connections.OutputConnection) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	conn, ok := m.conn(id)
	if ok {
		if c != nil && c.ID() != conn.ID() {
			return nil
		}
		err := m.onLogout(id, conn)
		if err != nil {
			return err
		}
	}
	m.Identities.Delete(id)
	return nil
}

//Verify if user given by id is useing given connection.
//Return verify result and any error raised.
func (m *Map) Verify(id string, conn connections.OutputConnection) (bool, error) {
	conn, ok := m.conn(id)
	if ok == false {
		return false, nil
	}
	return conn.ID() == id, nil
}

//SendByID send message to given user.
//Return any error if raised.
func (m *Map) SendByID(id string, msg []byte) error {
	conn, ok := m.conn(id)
	if ok == false {
		return nil
	}
	return conn.Send(msg)
}

//OnLogout return user logout callback function.
func (m *Map) OnLogout() func(id string, conn connections.OutputConnection) error {
	return m.onLogout
}

//SetOnLogout set user logout callback function.
func (m *Map) SetOnLogout(f func(id string, conn connections.OutputConnection) error) {
	m.onLogout = f
}

// NewMap create new map  identifier
// Logout callback of identifier will be gererated by  GenerateDefaultMapOnLogout.
func NewMap() *Map {
	m := &Map{
		Identities: sync.Map{},
	}
	m.onLogout = GenerateDefaultMapOnLogout(m)
	return m
}
