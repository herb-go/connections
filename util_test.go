package connections

import (
	"time"
)

func readByteChan(c chan []byte) ([]byte, bool) {
	select {
	case bs, more := <-c:
		return bs, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}

func readMessageChan(c chan *Message) (*Message, bool) {
	select {
	case v, more := <-c:
		return v, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}

func readConnChan(c chan OutputConnection) (OutputConnection, bool) {
	select {
	case v, more := <-c:
		return v, more
	case <-time.After(time.Millisecond):
		return nil, false
	}
}
