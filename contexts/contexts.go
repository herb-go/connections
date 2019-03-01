package contexts

import (
	"sync"

	"github.com/herb-go/connections"
)

type ConnContext struct {
	connections.OutputConnection
	Data sync.Map
	Lock sync.RWMutex
}

func NewConnContext() *ConnContext {
	return &ConnContext{}
}

type Contexts struct {
	connections.EmptyConsumer
	conns sync.Map
}

func (c *Contexts) OnClose(conn connections.OutputConnection) {
	id := conn.ID()
	c.conns.Delete(id)

}
func (c *Contexts) OnOpen(conn connections.OutputConnection) {
	id := conn.ID()
	context := NewConnContext()
	context.OutputConnection = conn
	c.conns.Store(id, context)
}

func (c *Contexts) Context(id string) *ConnContext {
	v, ok := c.conns.Load(id)
	if ok == false {
		return nil
	}
	return v.(*ConnContext)
}
