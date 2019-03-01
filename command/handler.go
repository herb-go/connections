package command

import (
	"github.com/herb-go/connections"
)

type Handler func(conn connections.OutputConnection, cmd Command) error

type Handlers map[string]Handler

func (h Handlers) WrapError(conn connections.OutputConnection, err error) *connections.Error {
	if err == nil {
		return nil
	}
	return &connections.Error{
		Conn:  conn,
		Error: err,
	}
}
func (h Handlers) Add(commandType string, handler Handler) {
	h[commandType] = handler
}
func (h Handlers) Exec(msg *connections.Message) (Command, bool, *connections.Error) {
	cmd := New()
	err := cmd.Decode(msg.Message)
	if err != nil {
		return nil, false, h.WrapError(msg.Conn, err)
	}

	handler, ok := h[cmd.Type()]
	if ok == false {
		return cmd, false, nil
	}
	err = handler(msg.Conn, cmd)
	if err != nil {
		return cmd, true, h.WrapError(msg.Conn, err)
	}
	return cmd, true, nil
}

func NewHandlers() Handlers {
	h := Handlers(map[string]Handler{})
	return h
}
