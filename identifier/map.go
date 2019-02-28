package identifier

import (
	"sync"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

var GenerateDefaultMapOnLogout = func(m *Map) func(id string, conn connections.ConnectionOutput) error {
	return func(id string, conn connections.ConnectionOutput) error {
		conn.Close()
		return nil
	}
}

type Map struct {
	Identities sync.Map
	lock       sync.Mutex
	onLogout   func(id string, conn connections.ConnectionOutput) error
}

func (m *Map) conn(id string) (connections.ConnectionOutput, bool) {
	data, ok := m.Identities.Load(id)
	if ok == false {
		return nil, false
	}
	conn, ok := data.(*connections.Conn)
	return conn, ok
}

func (m *Map) Login(id string, c connections.ConnectionOutput) error {
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
func (m *Map) Logout(id string, c connections.ConnectionOutput) error {
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
func (m *Map) Verify(id string, conn connections.ConnectionOutput) (bool, error) {
	conn, ok := m.conn(id)
	if ok == false {
		return false, nil
	}
	return conn.ID() == id, nil
}
func (m *Map) SendByID(id string, msg []byte) error {
	conn, ok := m.conn(id)
	if ok == false {
		return nil
	}
	return conn.Send(msg)
}
func (m *Map) OnLogout() func(id string, conn connections.ConnectionOutput) error {
	return m.onLogout
}
func (m *Map) SetOnLogout(f func(id string, conn connections.ConnectionOutput) error) {
	m.onLogout = f
}

func NewMap() *Map {

	m := &Map{
		Identities: sync.Map{},
	}
	m.onLogout = GenerateDefaultMapOnLogout(m)
	return m
}
