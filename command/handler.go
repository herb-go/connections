package command

import (
	"github.com/herb-go/connections"
)

// Handler command handler type.
type Handler func(conn connections.OutputConnection, cmd Command) error

var SeparatedCommandUnmarshaler = func(msg []byte) (Command, error) {
	c := New()
	err := c.Decode(msg)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//Handlers command handlers manager.
type Handlers struct {
	Handlers    map[string]Handler
	Unmarshaler func([]byte) (Command, error)
}

//WrapError wrap connection and error to connections.Error.
//Return connections.Error wrapped or nil if no error.
func (h *Handlers) WrapError(conn connections.OutputConnection, err error) *connections.Error {
	if err == nil {
		return nil
	}
	return &connections.Error{
		Conn:  conn,
		Error: err,
	}
}

// Register handler to given command type
func (h *Handlers) Register(commandType string, handler Handler) {
	h.Handlers[commandType] = handler
}

//Exec exec connection message
//Return convet decoded command ,whether handler for command type exists,and any connections error if raised.
func (h *Handlers) Exec(msg *connections.Message) (Command, bool, *connections.Error) {
	cmd, err := h.Unmarshaler(msg.Message)
	if err != nil {
		return nil, false, h.WrapError(msg.Conn, err)
	}

	handler, ok := h.Handlers[cmd.Type()]
	if ok == false {
		return cmd, false, nil
	}
	err = handler(msg.Conn, cmd)
	if err != nil {
		return cmd, true, h.WrapError(msg.Conn, err)
	}
	return cmd, true, nil
}

// NewHandlers create new handlers
func NewHandlers() *Handlers {
	h := &Handlers{
		Handlers:    map[string]Handler{},
		Unmarshaler: SeparatedCommandUnmarshaler,
	}
	return h
}
