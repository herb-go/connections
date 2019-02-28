package room

import (
	"container/list"
	"sync"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

type Location struct {
	list  *list.List
	lock  sync.Mutex
	conn  connections.ConnectionOutput
	rooms Joinable
}

func NewLocation(conn connections.ConnectionOutput, rooms Joinable) *Location {
	return &Location{
		list:  list.New(),
		conn:  conn,
		rooms: rooms,
	}
}

func (l *Location) Join(roomid string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.rooms.Join(roomid, l.conn)
	l.list.PushBack(roomid)
}
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
		rooms[i] = e.Value.(string)
		e = e.Next()
		i++
	}
	return rooms
}
