package room

import (
	"time"

	"github.com/herb-go/connections"
)

func readBytesChan(c chan []byte) ([]byte, bool) {
	select {
	case bs, more := <-c:
		return bs, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}

func readMessageChan(c chan *connections.Message) (*connections.Message, bool) {
	select {
	case v, more := <-c:
		return v, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}

func readErrorChan(c chan *connections.Error) (*connections.Error, bool) {
	select {
	case v, more := <-c:
		return v, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}

func readConnChan(c chan connections.OutputConnection) (connections.OutputConnection, bool) {
	select {
	case v, more := <-c:
		return v, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}
