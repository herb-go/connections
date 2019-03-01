package identifier

import "github.com/herb-go/connections"

type Identifier interface {
	Login(id string, conn connections.OutputConnection) error
	Logout(id string, conn connections.OutputConnection) error
	Verify(id string, conn connections.OutputConnection) (bool, error)
	SendByID(id string, msg []byte) error
	OnLogout() func(id string, conn connections.OutputConnection) error
	SetOnLogout(func(id string, conn connections.OutputConnection) error)
}
