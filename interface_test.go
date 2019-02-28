package connections

import "testing"

func TestInterface(t *testing.T) {
	var c ConnectionOutput
	c = New()
	id := c.ID()
	if id != "" {
		t.Error(id)
	}
	var c2 ConnectionsInput
	c2 = NewGateway()
	if c2 == nil {
		t.Error(c2)
	}
	var c3 ConnectionsOutput
	c3 = NewGateway()
	if c3 == nil {
		t.Error(c3)
	}
	var consumer ConnectionsConsumer
	consumer = EmptyConsumer{}
	if consumer == nil {
		t.Error(consumer)
	}
}
