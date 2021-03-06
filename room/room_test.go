package room

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/herb-go/connections"
)

var errConnectionSend = errors.New("connection send")

type errorChanConnection struct {
	*connections.ChanConnection
}

func (c *errorChanConnection) Send(msg []byte) error {
	return errConnectionSend
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
func TestRoom(t *testing.T) {
	var testroomid1 = "room1"
	var testroomid2 = "room2"
	var testroomidnotexist = "roomnotexist"
	var testmsg = []byte("test")
	var testmsg2 = []byte("test")
	var rooms = NewRooms()
	var g = connections.NewGateway()
	var tc = newTestConsumer()
	go func() {
		connections.Consume(g, tc)
	}()
	var chanconn1 = connections.NewChanConnection()
	var chanconn2 = connections.NewChanConnection()
	conn1, err := g.Register(chanconn1)
	if err != nil {
		t.Fatal(chanconn1)
	}
	conn2, err := g.Register(chanconn2)
	if err != nil {
		t.Fatal(chanconn2)
	}

	time.Sleep(time.Millisecond)
	location1 := NewLocation(conn1, rooms)
	location2 := NewLocation(conn2, rooms)
	location1.Join(testroomid1)
	location2.Join(testroomid2)

	membernotexsit := rooms.Members(testroomidnotexist)
	if len(membernotexsit) != 0 {
		t.Fatal(membernotexsit)
	}
	rooms.Broadcast(testroomid1, testmsg, nil)
	bs, _ := readBytesChan(chanconn1.ClientMessagesChan())
	if bytes.Compare(bs, testmsg) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bs != nil {
		t.Fatal(bs)
	}
	rooms.Broadcast(testroomid2, testmsg2, nil)
	bs, _ = readBytesChan(chanconn1.ClientMessagesChan())
	if bs != nil {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bytes.Compare(bs, testmsg2) != 0 {
		t.Fatal(bs)
	}
	location2.Leave(testroomid2)
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bs != nil {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bs != nil {
		t.Fatal(bs)
	}
	location2.Join(testroomid1)
	rooms.Broadcast(testroomid1, testmsg, nil)
	bs, _ = readBytesChan(chanconn1.ClientMessagesChan())
	if bytes.Compare(bs, testmsg) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bytes.Compare(bs, testmsg) != 0 {
		t.Fatal(bs)
	}
	testroom1members := rooms.Members(testroomid1)
	if len(testroom1members) != 2 {
		t.Fatal(testroom1members)
	}
	location1.Join(testroomid2)
	rooms.Broadcast(testroomid1, testmsg, nil)
	rooms.Broadcast(testroomid2, testmsg2, nil)
	bs, _ = readBytesChan(chanconn1.ClientMessagesChan())
	if bytes.Compare(bs, testmsg) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn1.ClientMessagesChan())
	if bytes.Compare(bs, testmsg2) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bytes.Compare(bs, testmsg2) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bs != nil {
		t.Fatal(bs)
	}

	memeber1rooms := location1.Rooms()
	if len(memeber1rooms) != 2 {
		t.Fatal(memeber1rooms)
	}
	memeber2rooms := location2.Rooms()
	if len(memeber2rooms) != 1 {
		t.Fatal(memeber1rooms)
	}
	location1.LeaveAll()
	rooms.Broadcast(testroomid1, testmsg, nil)
	rooms.Broadcast(testroomid2, testmsg2, nil)
	bs, _ = readBytesChan(chanconn1.ClientMessagesChan())
	if bytes.Compare(bs, testmsg) != 0 {
		t.Fatal(bs)
	}
	bs, _ = readBytesChan(chanconn2.ClientMessagesChan())
	if bytes.Compare(bs, testmsg2) != 0 {
		t.Fatal(bs)
	}

	errconnection := &errorChanConnection{
		ChanConnection: connections.NewChanConnection(),
	}
	connerr, err := g.Register(errconnection)
	if err != nil {
		t.Fatal(errconnection)
	}
	rooms.Join(testroomid2, connerr)
	var rerr *BroadcastError
	rooms.Broadcast(testroomid2, testmsg2, func(err error) {
		rerr = err.(*BroadcastError)
	})
	if rerr == nil || rerr.Room.ID != testroomid2 || rerr.Conn != connerr {
		t.Fatal(rerr)
	}
	r, ok := rooms.Rooms.Load(testroomid1)
	if ok == false {
		t.Fatal(ok)
	}
	ok = r.(*Room).Join(conn2)
	if ok {
		t.Fatal(ok)
	}
	ok = r.(*Room).Leave(conn2)
	if !ok {
		t.Fatal(ok)
	}
	ok = r.(*Room).Join(conn2)
	if !ok {
		t.Fatal(ok)
	}

	ok = r.(*Room).Leave(connerr)
	if ok {
		t.Fatal(ok)
	}
}
