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
var errID = errors.New("ID")
var errIDGenerator = func() (string, error) {
	return "error", errID
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
	time.Sleep(time.Millisecond)
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
func TestErrConnID(t *testing.T) {
	var conn *Conn
	var err error
	g := NewGateway()
	g.IDGenerator = errIDGenerator
	chanconn := NewChanConnection()
	conn, err = g.Register(chanconn)
	if conn != nil {
		t.Fatal(conn)
	}
	if err != errID {
		t.Fatal(err)
	}
}
func TestGateway(t *testing.T) {
	var testmsg = []byte("test")
	var testbackmsg = []byte("testback")
	var testerror = errors.New("test")
	var conn *Conn
	var err error
	var connIDNotExists = "notexists"
	g := NewGateway()
	g.ID = "test"
	conn = g.Conn(connIDNotExists)
	if conn != nil {
		t.Fatal(conn)
	}
	err = g.Send(connIDNotExists, testbackmsg)
	if err != nil {
		t.Fatal(err)
	}
	err = g.Close(connIDNotExists)
	if err != nil {
		t.Fatal(err)
	}
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
	conns := g.ListConn()
	if len(conns) != 1 {
		t.Fatal(conns)
	}
	err = g.Close(conn.ID())
	if err != nil {
		panic(err)
	}
	conns = g.ListConn()

	if len(conns) != 0 {
		t.Fatal(conns)
	}
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
	chanconn = NewChanConnection()
	conn, err = g.Register(chanconn)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if chanconn.Closed == true {
		t.Fatal(chanconn.Closed)
	}
	g.Stop()
	g.Stop()
	_, more = readConnChan(g.OnOpenEventsChan())
	if more != true {
		t.Fatal(more)
	}
	time.Sleep(time.Millisecond)
	if chanconn.Closed == false {
		t.Fatal(chanconn.Closed)
	}
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
