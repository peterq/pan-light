package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
	"math/rand"
	"sync"
)

type roleType interface {
	roleName() string
}

type roleHost struct {
	name    string
	session *realtime.Session
	slaves  map[string]*roleSlave
}

func (*roleHost) roleName() string {
	return "host"
}

type roleUser struct {
	session *realtime.Session

	waitState *waitState
}

func (*roleUser) roleName() string {
	return "user"
}

func (user *roleUser) requestTicket() (data gson, err error) {
	if user.waitState != nil {
		data = gson{
			"order":  user.waitState.order,
			"ticket": user.waitState.ticket,
		}
		return
	} else {
		manager.waitSessionMapLock.Lock()
		defer manager.waitSessionMapLock.Lock()
		manager.lastDistributedOrder++
		w := &waitState{
			ticket:  randomStr(32),
			order:   manager.lastDistributedOrder,
			session: user.session,
		}
		manager.waitSessionMap[manager.lastDistributedOrder] = w
		user.waitState = w
		server.RoomByName("room.all.host").Broadcast("wait.user.new", nil)
	}
	return
}

type roleSlave struct {
	name        string
	host        *roleHost
	session     *realtime.Session
	userSession *realtime.Session
	lock        sync.Mutex
}

func (*roleSlave) roleName() string {
	return "slave"
}

func randomStr(lenght int) string {
	arr := make([]byte, lenght)
	src := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	for i := 0; i < lenght; i++ {
		arr[i] = byte(src[rand.Intn(len(src))])
	}
	return string(arr)
}
