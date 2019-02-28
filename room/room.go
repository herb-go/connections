package room

import (
	"sync"

	"container/list"

	"github.com/jarlyyn/herb-go-experimental/connections"
)

type BroadcastError struct {
	Error error
	Conn  connections.ConnectionOutput
	Room  *Room
}
type Room struct {
	ID    string
	Lock  sync.Mutex
	Conns *list.List
}

func (r *Room) Members() []connections.ConnectionOutput {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	conns := make([]connections.ConnectionOutput, r.Conns.Len())
	e := r.Conns.Front()
	i := 0
	for {
		if e == nil {
			break
		}
		conns[i] = e.Value.(connections.ConnectionOutput)
		e = e.Next()
		i++
	}
	return conns
}
func (r *Room) Join(conn connections.ConnectionOutput) bool {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	newid := conn.ID()
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.ConnectionOutput)
		if c != nil && c.ID() == newid {
			return false
		}
		e = e.Next()
	}
	r.Conns.PushBack(conn)
	return true
}

func (r *Room) Leave(conn connections.ConnectionOutput) bool {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	newid := conn.ID()
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.ConnectionOutput)
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
		c := e.Value.(connections.ConnectionOutput)
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

func (r *Rooms) Members(roomid string) []connections.ConnectionOutput {
	v, ok := r.Rooms.Load(roomid)
	if ok == false || v == nil {
		return []connections.ConnectionOutput{}
	}
	return v.(*Room).Members()
}
func (r *Rooms) Join(roomid string, conn connections.ConnectionOutput) {
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

func (r *Rooms) Leave(roomid string, conn connections.ConnectionOutput) {
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
	Join(roomid string, conn connections.ConnectionOutput)
	Leave(roomid string, conn connections.ConnectionOutput)
}
