package connections

import (
	"bytes"
	"testing"
)

func TestChanConnection(t *testing.T) {
	var testmsg = []byte("test")
	var testmsgback = []byte("testback")
	var err error
	conn := NewChanConnection()
	client := conn.Client()
	err = client.Send(testmsg)
	if err != nil {
		t.Fatal(err)
	}
	m, ok := readBytesChan(conn.MessagesChan())
	if ok == false {
		t.Fatal(ok)
	}
	if bytes.Compare(m, testmsg) != 0 {
		t.Fatal(m)
	}
	err = conn.Send(testmsgback)
	if err != nil {
		t.Fatal(err)
	}
	m, ok = readBytesChan(client.MessagesChan())
	if ok == false {
		t.Fatal(ok)
	}
	if bytes.Compare(m, testmsgback) != 0 {
		t.Fatal(m)
	}
	if conn.Closed == true {
		t.Fatal(conn.Closed)
	}
	err = client.Close()
	if err != nil {
		t.Fatal(err)
	}
	if conn.Closed != true {
		t.Fatal(conn.Closed)
	}

}
