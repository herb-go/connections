package websocket

import (
	"testing"

	"github.com/herb-go/connections"
)

func TestInterface(t *testing.T) {
	var c connections.Conn = New()
	c.Close()
}
