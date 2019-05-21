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
	if user.waitState == nil {
		manager.waitSessionMapLock.Lock()
		defer manager.waitSessionMapLock.Unlock()
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
	data = gson{
		"order":  user.waitState.order,
		"ticket": user.waitState.ticket,
	}
	return
}

type slaveState string

const (
	slaveStateWait     slaveState = "wait"
	slaveStateStarting slaveState = "starting"
	slaveStateRuning   slaveState = "running"
)

type roleSlave struct {
	name          string            // slave 名称, 需要已host名称为前缀, 用来鉴权
	host          *roleHost         // 指向host
	session       *realtime.Session // slave 进程链接的session
	userWaitState *waitState        // 用户排队票据
	state         slaveState
	lock          sync.Mutex
}

func (*roleSlave) roleName() string {
	return "slave"
}

func randomStr(length int) string {
	arr := make([]byte, length)
	src := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	for i := 0; i < length; i++ {
		arr[i] = byte(src[rand.Intn(len(src))])
	}
	return string(arr)
}
