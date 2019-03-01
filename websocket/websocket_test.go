package websocket

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/herb-go/connections"
)

func TestInterface(t *testing.T) {
	var c connections.RawConnection = New()
	c.Close()
}

func TestMethods(t *testing.T) {
	var conn *Conn
	var testmsg = []byte("test")
	var testmsgback = []byte("testback")
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.Handle("/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		conn, err = Upgrade(w, r, MsgTypeText)
		if err != nil {
			t.Fatal(err)
		}
	}))
	wsurl := strings.Replace(server.URL, "http://", "ws://", 1) + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		t.Fatal(err)
	}
	ws.WriteMessage(MsgTypeText, testmsg)
	msg := <-conn.MessagesChan()
	if bytes.Compare(msg, testmsg) != 0 {
		t.Fatal(msg)
	}
	ws.WriteMessage(MsgTypeBinary, testmsg)
	cerr := <-conn.ErrorsChan()
	if cerr != ErrMsgTypeNotMatch {
		t.Fatal(cerr)
	}
	err = conn.Send(testmsgback)
	if err != nil {
		t.Fatal(err)
	}
	ty, msg, err := ws.ReadMessage()
	if ty != MsgTypeText {
		t.Fatal(t)
	}
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(msg, testmsgback) != 0 {
		t.Fatal(msg)
	}
	addr := ws.LocalAddr()
	if addr.String() != conn.RemoteAddr().String() {
		t.Fatal(addr)
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, closed := <-conn.C()
	if closed != false {
		t.Fatal(closed)
	}
	err = conn.send(testmsgback)
	if err != nil {
		t.Fatal(err)
	}
	ws, _, err = websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		t.Fatal(err)
	}
	err = ws.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, closed = <-conn.C()
	if closed != false {
		t.Fatal(closed)
	}

}
