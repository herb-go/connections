package identifier

import "testing"

func TestInterface(t *testing.T) {
	var m Identifier
	m = NewMap()
	f := m.OnLogout()
	if f == nil {
		t.Error(&f)
	}
}
