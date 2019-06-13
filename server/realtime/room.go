package realtime

import "sync"

type sessionIdSlice []SessionId

func (slice sessionIdSlice) index(id SessionId) int {
	for idx, idInSlice := range slice {
		if idInSlice == id {
			return idx
		}
	}
	return -1
}

type Room struct {
	name    string
	server  *Server
	members sessionIdSlice
	lock    sync.RWMutex
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Join(id SessionId) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.members.index(id) > -1 {
		return
	}
	r.members = append(r.members, id)

	ss, _ := r.server.SessionById(id)
	ss.lock.Lock()
	ss.rooms = append(ss.rooms, r)
	ss.lock.Unlock()
	go r.Broadcast("room.member.join", id)
}

func (r *Room) Count() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.members)
}

func (r *Room) Broadcast(event string, payload interface{}, expect ...SessionId) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	for _, id := range r.members {
		if sessionIdSlice(expect).index(id) < 0 {
			ss, ok := r.server.SessionById(id)
			if !ok {
				go r.Remove(id)
				continue
			}
			ss.Emit(event, payload, r.name)
		}
	}
}

func (r *Room) Remove(id SessionId) {
	func() {
		ss, _ := r.server.SessionById(id)
		if ss == nil {
			return
		}
		ss.lock.Lock()
		defer ss.lock.Unlock()
		rooms := make([]*Room, 0, cap(ss.rooms))
		for _, room := range ss.rooms {
			if room != r {
				rooms = append(rooms, room)
			}
		}
		ss.rooms = rooms
	}()
	r.lock.Lock()
	defer r.lock.Unlock()
	idx := r.members.index(id)
	if idx == -1 {
		return
	}
	go r.Broadcast("room.member.remove", id)
	go func() {
		if ss, ok := r.server.SessionById(id); ok {
			ss.Emit("room.leave", r.name)
		}
	}()
	if idx == len(r.members)-1 {
		r.members = r.members[:idx]
		return
	}
	if idx == 0 {
		r.members = r.members[1:]
		return
	}
	r.members = append(r.members[:idx], r.members[idx+1:]...)
}

func (r *Room) Members() []SessionId {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.members
}
