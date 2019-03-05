package command

import (
	"errors"
	"testing"

	"github.com/herb-go/connections"
)

var errDecodeError = errors.New("deode error")
var errUnmarshaler = func(msg []byte) (Command, error) {
	return nil, errDecodeError
}

func TestWrapError(t *testing.T) {
	hs := NewHandlers()
	var conn *connections.Conn
	werr := hs.WrapError(conn, nil)
	if werr != nil {
		t.Fatal(werr)
	}
	hs.Unmarshaler = errUnmarshaler
	msg := &connections.Message{
		Conn:    nil,
		Message: []byte{},
	}
	_, _, werr = hs.Exec(msg)
	if werr == nil || werr.Error != errDecodeError {
		t.Fatal(werr)
	}
}
func TestHandler(t *testing.T) {
	var dummyconn = connections.NewDummyConnection()
	g := connections.NewGateway()
	go func() {
		connections.Consume(g, connections.EmptyConsumer{})
	}()
	var conn, err = g.Register(dummyconn)
	if err != nil {
		t.Fatal(err)
	}
	var lastcmds Command
	var lastcalled string
	var type1 = "type1"
	var typeerror = "typeerror"
	var typeNotExist = "notexists"
	var testerr = errors.New("test error")
	hs := NewHandlers()
	hs.Register(type1, func(conn connections.OutputConnection, cmd Command) error {
		lastcmds = cmd
		lastcalled = type1
		return nil
	})
	hs.Register(typeerror, func(conn connections.OutputConnection, cmd Command) error {
		return testerr
	})
	m := &connections.Message{
		Conn:    conn,
		Message: []byte(type1),
	}
	c, ok, cerr := hs.Exec(m)
	if ok != true {
		t.Fatal(ok)
	}
	if cerr != nil {
		t.Fatal(cerr)
	}
	if c.Type() != type1 {
		t.Fatal(c.Type())
	}
	if lastcmds != c {
		t.Fatal(lastcmds)
	}
	if lastcalled != type1 {
		t.Fatal(type1)
	}

	m = &connections.Message{
		Conn:    conn,
		Message: []byte(typeNotExist),
	}
	c, ok, cerr = hs.Exec(m)
	if c.Type() != typeNotExist {
		t.Fatal(c)
	}
	if ok != false {
		t.Fatal(ok)
	}
	if cerr != nil {
		t.Fatal(cerr)
	}

	m = &connections.Message{
		Conn:    conn,
		Message: []byte(typeerror),
	}
	c, ok, cerr = hs.Exec(m)
	if c.Type() != typeerror {
		t.Fatal(c)
	}
	if ok != true {
		t.Fatal(ok)
	}
	if cerr == nil || cerr.Error != testerr {
		t.Fatal(cerr)
	}
}
