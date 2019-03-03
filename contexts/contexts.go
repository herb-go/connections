package contexts

import (
	"sync"

	"github.com/herb-go/connections"
)

// ConnContext context of connection
type ConnContext struct {
	connections.OutputConnection
	Data sync.Map
	Lock sync.RWMutex
}

//NewConnContext create new connection context.
//return connection context created
func NewConnContext() *ConnContext {
	return &ConnContext{}
}

//Contexts connections contexts.
//You should extend this struct to your own connection consumer.
type Contexts struct {
	connections.EmptyConsumer
	conns sync.Map
}

//OnClose called when connection closed.
//You should call this method if you overwrited it.
func (c *Contexts) OnClose(conn connections.OutputConnection) {
	id := conn.ID()
	c.conns.Delete(id)

}

//OnOpen called when connection open.
//You should call this method if you overwrited it.
func (c *Contexts) OnOpen(conn connections.OutputConnection) {
	id := conn.ID()
	context := NewConnContext()
	context.OutputConnection = conn
	c.conns.Store(id, context)
}

//Context get connection context by given connection id.
//Return context or nil if connection not found.
func (c *Contexts) Context(id string) *ConnContext {
	v, ok := c.conns.Load(id)
	if ok == false {
		return nil
	}
	return v.(*ConnContext)
}

//New create new Contexts.
func New() *Contexts {
	return &Contexts{}
}
