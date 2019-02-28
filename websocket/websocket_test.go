package websocket

import (
	"testing"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

func TestInterface(t *testing.T) {
	var c connections.Conn = New()
	c.Close()
}
