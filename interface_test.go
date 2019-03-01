package connections

import "testing"

func TestInterface(t *testing.T) {
	var c OutputConnection
	c = New()
	id := c.ID()
	if id != "" {
		t.Error(id)
	}
	var c2 InputService
	c2 = NewGateway()
	if c2 == nil {
		t.Error(c2)
	}
	var c3 OutputService
	c3 = NewGateway()
	if c3 == nil {
		t.Error(c3)
	}
	var consumer Consumer
	consumer = EmptyConsumer{}
	if consumer == nil {
		t.Error(consumer)
	}
	var rc RawConnection
	rc = NewDummyConnection()
	if rc == nil {
		t.Error(rc)
	}
}
