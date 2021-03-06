package connections

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

type testConsumer struct {
	EmptyConsumer
	lastMessage *Message
	lastError   *Error
	lastOpen    OutputConnection
	lastClose   OutputConnection
	stoped      bool
}

//OnMessage called when connection message received.
func (c *testConsumer) OnMessage(m *Message) {
	c.lastMessage = m
}

//OnError called when onconnection error raised.
func (c *testConsumer) OnError(e *Error) {
	c.lastError = e
}

//OnClose called when connection closed.
func (c *testConsumer) OnClose(oc OutputConnection) {
	c.lastClose = oc
}

//OnOpen called when connection open.
func (c *testConsumer) OnOpen(oc OutputConnection) {
	c.lastOpen = oc
}

// Stop stop consumer
func (c *testConsumer) Stop() {
	c.stoped = true
}

func TestConsumer(t *testing.T) {
	var testmsg = []byte("test")
	var testerror = errors.New("test error")
	var tc = &testConsumer{}
	var g = NewGateway()
	var chanconn = NewChanConnection()
	go func() {
		Consume(g, tc)
	}()
	time.Sleep(time.Millisecond)
	conn, err := g.Register(chanconn)
	if err != nil {
		t.Fatal(err)
	}
	if conn == nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if tc.lastOpen != conn {
		t.Fatal(tc.lastOpen)
	}
	err = chanconn.ClientSend(testmsg)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)

	if tc.lastMessage == nil {
		t.Fatal(tc.lastMessage)
	}
	if tc.lastMessage.Conn != conn {
		t.Fatal(tc.lastMessage.Conn)
	}
	if bytes.Compare(tc.lastMessage.Message, testmsg) != 0 {
		t.Fatal(tc.lastMessage.Message)
	}

	chanconn.RaiseError(testerror)
	time.Sleep(time.Millisecond)
	if tc.lastError == nil {
		t.Error(tc.lastError)
	}
	if tc.lastError.Conn != conn {
		t.Error(tc.lastError.Conn)
	}
	if tc.lastError.Error != testerror {
		t.Error(tc.lastError.Error)
	}
	chanconn.Close()
	time.Sleep(time.Millisecond)
	if tc.lastClose != conn {
		t.Fatal(tc.lastClose)
	}
	if tc.stoped != false {
		t.Fatal(tc.stoped)
	}
	g.Stop()
	time.Sleep(time.Millisecond)
	if tc.stoped != true {
		t.Fatal(tc.stoped)
	}
}
