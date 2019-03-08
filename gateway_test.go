package connections

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"
	"time"
)

var duplicatedIDGenerator = func() (string, error) {
	return "duplicated id", nil
}

func TestErrConnIDDuplicated(t *testing.T) {
	var testmsg = []byte("test")
	var conn *Conn
	var err error
	g := NewGateway()
	g.IDGenerator = duplicatedIDGenerator
	chanconn1 := NewChanConnection()
	chanconn2 := NewChanConnection()
	go func() {
		conn, err = g.Register(chanconn1)
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
	time.Sleep(time.Microsecond)
	go func() {
		_, err = g.Register(chanconn2)
		if err != ErrConnIDDuplicated {
			t.Fatal(err)
		}
		c, more := readConnChan(g.OnOpenEventsChan())
		if more == true {
			t.Fatal(more)
		}
		if c != nil {
			t.Fatal(conn)
		}

	}()
	time.Sleep(time.Microsecond)
	chanconn1.ClientSend(testmsg)
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
	if chanconn2.Closed != true {
		t.Fatal(chanconn2.Closed)
	}
}

func TestGateway(t *testing.T) {
	var testmsg = []byte("test")
	var testbackmsg = []byte("testback")
	var testerror = errors.New("test")
	var conn *Conn
	var err error
	g := NewGateway()
	g.ID = "test"
	time.Sleep(time.Microsecond)
	chanconn := NewChanConnection()
	go func() {
		conn, err = g.Register(chanconn)
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
	time.Sleep(time.Second)
	chanconn.ClientSend(testmsg)
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
	bs, ok := readBytesChan(chanconn.ClientMessagesChan())
	if ok != true {
		t.Fatal(ok)
	}
	if bytes.Compare(bs, testbackmsg) != 0 {
		t.Fatal(bs)
	}
	chanconn.RaiseError(testerror)

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

	chanconn.Addr = &net.IPAddr{}

	if conn.RemoteAddr() != chanconn.Addr {
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
	if chanconn.Closed != true {
		t.Fatal(chanconn.Closed)
	}
	g.Stop()
	time.Sleep(time.Millisecond)
	chanconn = NewChanConnection()
	go func() {
		conn, err = g.Register(chanconn)
		if err != nil {
			t.Fatal(err)
		}
		_, more := readConnChan(g.OnOpenEventsChan())
		if more != false {
			t.Fatal(more)
		}

	}()
	err = chanconn.ClientSend(testmsg)
	if err != nil {
		t.Fatal(err)
	}
	_, more = readMessageChan(g.MessagesChan())
	if more != false {
		t.Fatal(more)
	}
	chanconn.RaiseError(errors.New("error"))
	_, more = readErrorChan(g.ErrorsChan())
	if more != false {
		t.Fatal(more)
	}

	time.Sleep(time.Millisecond)
	chanconn.Close()
	_, more = readConnChan(g.OnCloseEventsChan())
	if more != false {
		t.Fatal(more)
	}

}
