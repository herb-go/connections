package contexts

import (
	"sync"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

type ConnContext struct {
	connections.ConnectionOutput
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

func (c *Contexts) OnClose(conn connections.ConnectionOutput) {
	id := conn.ID()
	c.conns.Delete(id)

}
func (c *Contexts) OnOpen(conn connections.ConnectionOutput) {
	id := conn.ID()
	context := NewConnContext()
	context.ConnectionOutput = conn
	c.conns.Store(id, context)
}

func (c *Contexts) Context(id string) *ConnContext {
	v, ok := c.conns.Load(id)
	if ok == false {
		return nil
	}
	return v.(*ConnContext)
}
