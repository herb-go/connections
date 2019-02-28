package identifier

import "github.com/jarlyyn/herb-go-experimental/connections"

type Identifier interface {
	Login(id string, conn connections.ConnectionOutput) error
	Logout(id string, conn connections.ConnectionOutput) error
	Verify(id string, conn connections.ConnectionOutput) (bool, error)
	SendByID(id string, msg []byte) error
	OnLogout() func(id string, conn connections.ConnectionOutput) error
	SetOnLogout(func(id string, conn connections.ConnectionOutput) error)
}
