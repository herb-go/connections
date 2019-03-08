package identifier

import (
	"bytes"
	"testing"

	"github.com/herb-go/connections"
)

var mapOnLogoutWithoutClose = func(id string, conn connections.OutputConnection) error {
	return nil
}

type testConsumer struct {
	connections.EmptyConsumer
	MessagesChan chan *connections.Message
}

func (c *testConsumer) OnMessage(m *connections.Message) {
	c.MessagesChan <- m
}

func newTestConsumer() *testConsumer {
	return &testConsumer{
		MessagesChan: make(chan *connections.Message, 100),
	}
}
func TestInterface(t *testing.T) {
	var m Identifier
	m = New()
	f := m.OnLogout()
	if f == nil {
		t.Error(&f)
	}
}

func TestMap(t *testing.T) {
	var testmsg = []byte("testmsg")
	m := New()
	g := connections.NewGateway()
	chanconn := connections.NewChanConnection()
	chanconn2 := connections.NewChanConnection()
	chanconn3 := connections.NewChanConnection()

	uidtest := "test"
	uidnotexist := "testnotexist"
	var tc = newTestConsumer()
	go func() {
		connections.Consume(g, tc)
	}()
	err := m.SendByID(uidnotexist, testmsg)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := g.Register(chanconn)
	if err != nil {
		t.Fatal(err)
	}
	conn2, err := g.Register(chanconn2)
	if err != nil {
		t.Fatal(err)
	}
	conn3, err := g.Register(chanconn3)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := m.Verify(uidtest, conn)
	if err != nil {
		t.Fatal(err)
	}
	if ok == true {
		t.Fatal(ok)
	}
	err = m.Login(uidtest, conn)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn)
	if err != nil {
		t.Fatal(err)
	}
	if ok == false {
		t.Fatal(ok)
	}
	err = m.SendByID(uidtest, testmsg)
	if err != nil {
		t.Fatal(err)
	}
	msg := <-chanconn.Output
	if bytes.Compare(msg, testmsg) != 0 {
		t.Fatal(m)
	}
	err = m.Login(uidtest, conn2)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn)
	if err != nil {
		t.Fatal(err)
	}
	if ok == true {
		t.Fatal(ok)
	}
	ok, err = m.Verify(uidtest, conn2)
	if err != nil {
		t.Fatal(err)
	}
	if ok == false {
		t.Fatal(ok)
	}
	err = m.Logout(uidtest, conn)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn2)
	if err != nil {
		t.Fatal(err)
	}
	if ok == false {
		t.Fatal(ok)
	}
	err = m.Logout(uidtest, conn2)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn2)
	if err != nil {
		t.Fatal(err)
	}
	if ok == true {
		t.Fatal(ok)
	}
	m.SetOnLogout(mapOnLogoutWithoutClose)
	err = m.Login(uidtest, conn3)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn3)
	if err != nil {
		t.Fatal(err)
	}
	if ok == false {
		t.Fatal(ok)
	}
	err = m.Logout(uidtest, nil)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = m.Verify(uidtest, conn3)
	if err != nil {
		t.Fatal(err)
	}
	if ok == true {
		t.Fatal(ok)
	}

	if !chanconn.Closed {
		t.Fatal(chanconn.Closed)
	}
	if !chanconn2.Closed {
		t.Fatal(chanconn2.Closed)
	}
	if chanconn3.Closed {
		t.Fatal(chanconn3.Closed)
	}
}
