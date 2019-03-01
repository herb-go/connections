package room

import (
	"sync"

	"container/list"

	"github.com/herb-go/connections"
)

type BroadcastError struct {
	Error error
	Conn  connections.OutputConnection
	Room  *Room
}
type Room struct {
	ID    string
	Lock  sync.Mutex
	Conns *list.List
}

func (r *Room) Members() []connections.OutputConnection {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	conns := make([]connections.OutputConnection, r.Conns.Len())
	e := r.Conns.Front()
	i := 0
	for {
		if e == nil {
			break
		}
		conns[i] = e.Value.(connections.OutputConnection)
		e = e.Next()
		i++
	}
	return conns
}
func (r *Room) Join(conn connections.OutputConnection) bool {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	newid := conn.ID()
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.OutputConnection)
		if c != nil && c.ID() == newid {
			return false
		}
		e = e.Next()
	}
	r.Conns.PushBack(conn)
	return true
}

func (r *Room) Leave(conn connections.OutputConnection) bool {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	newid := conn.ID()
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.OutputConnection)
		if c != nil && c.ID() == newid {
			r.Conns.Remove(e)
			return true
		}
		e = e.Next()
	}
	return false
}

func (r *Room) Broadcast(msg []byte) []*BroadcastError {
	errs := []*BroadcastError{}
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.OutputConnection)
		err := c.Send(msg)
		if err != nil {
			e := &BroadcastError{
				Error: err,
				Conn:  c,
				Room:  r,
			}
			errs = append(errs, e)
		}
		e = e.Next()
	}
	return errs
}
func NewRoom() *Room {
	return &Room{
		Conns: list.New(),
	}
}

type Rooms struct {
	Rooms  sync.Map
	Lock   sync.Mutex
	Errors chan *BroadcastError
}

func (r *Rooms) Members(roomid string) []connections.OutputConnection {
	v, ok := r.Rooms.Load(roomid)
	if ok == false || v == nil {
		return []connections.OutputConnection{}
	}
	return v.(*Room).Members()
}
func (r *Rooms) Join(roomid string, conn connections.OutputConnection) {
	var room *Room
	v, ok := r.Rooms.Load(roomid)
	if ok == false {
		r.Lock.Lock()
		room = NewRoom()
		room.ID = roomid
		v, _ = r.Rooms.LoadOrStore(roomid, room)
		r.Lock.Unlock()
	}
	room = v.(*Room)
	room.Join(conn)
	return
}

func (r *Rooms) Leave(roomid string, conn connections.OutputConnection) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	var room *Room
	v, ok := r.Rooms.Load(roomid)
	if ok == false {
		return
	}
	room = v.(*Room)
	ok = room.Leave(conn)
	if ok == false {
		return
	}
	if room.Conns.Len() == 0 {
		r.Rooms.Delete(roomid)
	}
	return
}

func (r *Rooms) Broadcast(roomid string, msg []byte) {
	var room *Room
	v, ok := r.Rooms.Load(roomid)
	if ok == false {
		return
	}
	room = v.(*Room)
	errs := room.Broadcast(msg)
	for i := range errs {
		r.Errors <- errs[i]
	}
	return
}
func NewRooms() *Rooms {
	return &Rooms{
		Errors: make(chan *BroadcastError),
	}
}

type Joinable interface {
	Join(roomid string, conn connections.OutputConnection)
	Leave(roomid string, conn connections.OutputConnection)
}
