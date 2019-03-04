package identifier

import "github.com/herb-go/connections"

// Identifier connections identifier interface.
type Identifier interface {
	//Login given connection as user by given id.
	//Return any error if raised.
	Login(id string, conn connections.OutputConnection) error
	//Logout logout user by given id if current user is given connection.
	//Always logout user if given connection is nil.
	//Return any error if raised.
	Logout(id string, conn connections.OutputConnection) error
	//Verify if user given by id is useing given connection.
	//Return verify result and any error raised.
	Verify(id string, conn connections.OutputConnection) (bool, error)
	//SendByID send message to given user.
	//Return any error if raised.
	SendByID(id string, msg []byte) error
	//OnLogout return user logout callback function.
	OnLogout() func(id string, conn connections.OutputConnection) error
	//SetOnLogout set user logout callback function.
	SetOnLogout(func(id string, conn connections.OutputConnection) error)
}
