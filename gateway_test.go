package connections

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

func TestGateway(t *testing.T) {
	var testmsg = []byte("test")
	var testbackmsg = []byte("testback")
	var testerror = errors.New("test")
	var conn *Conn
	var err error
	g := NewGateway()
	g.ID = "test"
	time.Sleep(time.Microsecond)
	dummyconn := NewDummyConnection()
	go func() {
		conn, err = g.Register(dummyconn)
		if err != nil {
			t.Fatal(err)
		}
		c, more := readConnChan(g.OnOpenEventsChan())
		if more != true {
			t.Fatal(more)
		}
		if c != conn {
			t.Fatal(conn)
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
	err = g.Send(conn.ID(), testbackmsg)
	if err != nil {
		t.Fatal(err)
	}
	bs, ok := readBytesChan(dummyconn.Output)
	if ok != true {
		t.Fatal(ok)
	}
	if bytes.Compare(bs, testbackmsg) != 0 {
		t.Fatal(bs)
	}
	dummyconn.RaiseError(testerror)

	e, ok := readErrorChan(g.ErrorsChan())
	if ok != true {
		t.Fatal(ok)
	}
	if e == nil {
		t.Fatal(e)
	}
	if e.Conn != conn {
		t.Fatal(e.Conn)
	}
	if e.Error != testerror {
		t.Fatal(e.Error)
	}

	if !strings.HasPrefix(conn.ID(), g.ID+"-") {
		t.Fatal(conn.ID())
	}

	dummyconn.Addr = &net.IPAddr{}

	if conn.RemoteAddr() != dummyconn.Addr {
		t.Fatal(conn.RemoteAddr())
	}
	g.Close(conn.ID())
	c, more := readConnChan(g.OnCloseEventsChan())
	if more != true {
		t.Fatal(more)
	}
	if c != conn {
		t.Fatal(conn)
	}

}
