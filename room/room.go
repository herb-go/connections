package room

import (
	"fmt"
	"sync"

	"container/list"

	"github.com/herb-go/connections"
)

// BroadcastError room broadcast error
type BroadcastError struct {
	//Err raw error.
	Err error
	//Conn connections in which error raised.
	Conn connections.OutputConnection
	//Room room in which error raised.
	Room *Room
}

func (e *BroadcastError) Error() string {
	return fmt.Sprintf("boradcast error %s", e.Err.Error())
}

//Room connection room in which all connections will receive broadcast.
type Room struct {
	ID    string
	Lock  sync.Mutex
	Conns *list.List
}

//Members list room connections.
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

//Join join connection to room.
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

//Leave leave connection from room.
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

//Broadcast broadcast message to all connection in room.
//Return BroadcastError if any error raised.
func (r *Room) Broadcast(msg []byte, errHandler func(error)) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	e := r.Conns.Front()
	for {
		if e == nil {
			break
		}
		c := e.Value.(connections.OutputConnection)
		go func() {
			err := c.Send(msg)
			if err != nil {
				e := &BroadcastError{
					Err:  err,
					Conn: c,
					Room: r,
				}
				if errHandler != nil {
					errHandler(e)
				}
			}
		}()
		e = e.Next()
	}
}

// NewRoom create new room.
func NewRoom() *Room {
	return &Room{
		Conns: list.New(),
	}
}

// Rooms rooms manager
type Rooms struct {
	Rooms sync.Map
	Lock  sync.Mutex
}

// Members list connections in room by given room id.
func (r *Rooms) Members(roomid string) []connections.OutputConnection {
	v, ok := r.Rooms.Load(roomid)
	if ok == false || v == nil {
		return []connections.OutputConnection{}
	}
	return v.(*Room).Members()
}

//Join join connection to room by given room id.
//Auto create room if not exists.
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

//Leave leave connection form room by give room id.
//Auto remove room if  root empty.
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

//Broadcast brodcat message to give room.
//BroadcastError will be sent to Error chan if any error raised.
func (r *Rooms) Broadcast(roomid string, msg []byte, errHandler func(err error)) {
	var room *Room
	v, ok := r.Rooms.Load(roomid)
	if ok == false {
		return
	}
	room = v.(*Room)
	room.Broadcast(msg, errHandler)
}

// NewRooms create new rooms manager.
func NewRooms() *Rooms {
	return &Rooms{}
}

//Joinable  interfacer for which can joined as room manager.
type Joinable interface {
	Join(roomid string, conn connections.OutputConnection)
	Leave(roomid string, conn connections.OutputConnection)
}
