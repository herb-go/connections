package room

import (
	"container/list"
	"sync"

	"github.com/herb-go/connections"
)

// Location room data of user.
type Location struct {
	list  *list.List
	lock  sync.Mutex
	conn  connections.OutputConnection
	rooms Joinable
}

//NewLocation create new loction with given connection and rooms.
func NewLocation(conn connections.OutputConnection, rooms Joinable) *Location {
	return &Location{
		list:  list.New(),
		conn:  conn,
		rooms: rooms,
	}
}

//Join join user to give room
func (l *Location) Join(roomid string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.rooms.Join(roomid, l.conn)
	l.list.PushBack(roomid)
}

//Leave leave from given room.
func (l *Location) Leave(roomid string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	e := l.list.Front()
	for {
		if e == nil {
			break
		}
		id := e.Value.(string)
		if id == roomid {
			l.list.Remove(e)
			l.rooms.Leave(roomid, l.conn)
		}
		e = e.Next()
	}
}

//LeaveAll Leave from all rooms.
func (l *Location) LeaveAll() {
	l.lock.Lock()
	defer l.lock.Unlock()
	e := l.list.Front()
	for {
		if e == nil {
			break
		}
		id := e.Value.(string)
		l.rooms.Leave(id, l.conn)
		l.list.Remove(e)
		e = e.Next()
	}
}

// Rooms list all rooms user joined.
func (l *Location) Rooms() []string {
	l.lock.Lock()
	defer l.lock.Unlock()
	rooms := []string{}
	var i = 0
	e := l.list.Front()
	for {
		if e == nil {
			break
		}
		rooms = append(rooms, e.Value.(string))
		e = e.Next()
		i++
	}
	return rooms
}
