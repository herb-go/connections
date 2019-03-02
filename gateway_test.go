package connections

import (
	"bytes"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	var testmsg = []byte("test")
	var conn *Conn
	var err error
	g := NewGateway()
	go func() {
		Consume(g, &EmptyConsumer{})
	}()
	time.Sleep(time.Microsecond)
	dummyconn := NewDummyConnection()
	go func() {
		conn, err = g.Register(dummyconn)
		if err != nil {
			t.Fatal(err)
		}
	}()
	dummyconn.ClientSend(testmsg)
	m, more := readMessageChan(g.MessagesChan())
	if more != true {
		t.Fatal(more)
	}
	if m == nil {
		t.Fatal(m)
	}
	if m.Conn != conn {
		t.Fatal(conn)
	}
	if bytes.Compare(m.Message, testmsg) != 0 {
		t.Fatal(m.Message)
	}
	conn2 := g.Conn(conn.ID())
	if conn2 != conn {
		t.Fatal(conn2)
	}
}
