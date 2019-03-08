package contexts

import (
	"testing"
	"time"

	"github.com/herb-go/connections"
)

func TestContexts(t *testing.T) {
	var testkey = "test"
	var testvalue = "testvalue"
	g := connections.NewGateway()
	chanconn := connections.NewChanConnection()
	var conn *connections.Conn
	var err error
	contexts := New()
	go func() {
		connections.Consume(g, contexts)
	}()
	go func() {
		conn, err = g.Register(chanconn)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Millisecond)
	c := contexts.Context(conn.ID())
	c.Data.Store(testkey, testvalue)
	c2 := contexts.Context(conn.ID())
	v, _ := c2.Data.Load(testkey)
	if v.(string) != testvalue {
		t.Fatal(v)
	}
	chanconn.Close()
	time.Sleep(time.Millisecond)
	c3 := contexts.Context(conn.ID())
	if c3 != nil {
		t.Fatal(c3)
	}
}
