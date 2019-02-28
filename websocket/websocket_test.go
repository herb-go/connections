package websocket

import (
	"testing"

	"connections"
)

func TestInterface(t *testing.T) {
	var c connections.Conn = New()
	c.Close()
}
