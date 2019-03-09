package message

import (
	"errors"
	"testing"
)

func TestMessage(t *testing.T) {
	var output = make(chan *Message, 10)
	var testerr = errors.New("error")
	var errhandler = func(*Message) error {
		return testerr
	}
	var testhandler = func(m *Message) error {
		output <- m
		return nil
	}
	var testhandlertype = "test"
	var errhandlertype = "err"
	var notexistshandlertype = "notexists"
	adapter := NewAdapter()
	adapter.Register(errhandlertype, errhandler)
	adapter.Register(testhandlertype, testhandler)
	testmsg := New()
	testmsg.Type = testhandlertype
	testmsg.Data = "test"
	ok, err := adapter.Exec(testmsg)
	if err != nil {
		t.Fatal(err)
	}
	if ok == false {
		t.Fatal(ok)
	}
	mo := <-output
	if mo != testmsg {
		t.Fatal(mo)
	}
	testmsgnotexist := New()
	testmsgnotexist.Type = notexistshandlertype
	testmsgnotexist.Data = "test"
	ok, err = adapter.Exec(testmsgnotexist)
	if err != nil {
		t.Fatal(err)
	}
	if ok == true {
		t.Fatal(ok)
	}

	testmsgerr := New()
	testmsgerr.Type = errhandlertype
	testmsgerr.Data = "test"
	ok, err = adapter.Exec(testmsgerr)
	if err != testerr {
		t.Fatal(err)
	}
}
