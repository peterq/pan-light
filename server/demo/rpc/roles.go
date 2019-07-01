package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
	"math/rand"
	"sync"
)

type roleType interface {
	roleName() string
	publicInfo() gson
}

type roleHost struct {
	name       string
	session    *realtime.Session
	wsAgentUrl string
	slaves     map[string]*roleSlave
}

func (*roleHost) publicInfo() gson {
	panic("implement me")
}

func (*roleHost) roleName() string {
	return "host"
}

type roleUser struct {
	session *realtime.Session

	nickname string // 随机分配花名
	avatar   string // 随机分配花名

	waitState *waitState
}

func (user *roleUser) publicInfo() gson {
	return gson{
		"nickname": user.nickname,
		"avatar":   user.avatar,
	}
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
	slaveStateRunning  slaveState = "running"
)

type roleSlave struct {
	name          string            // slave 名称, 需要已host名称为前缀, 用来鉴权
	host          *roleHost         // 指向host
	session       *realtime.Session // slave 进程链接的session
	userWaitState *waitState        // 用户排队票据
	startTime     int               // 启动时间
	endTime       int               // 结束时间
	state         slaveState
	lock          sync.Mutex
}

func (slave *roleSlave) publicInfo() gson {
	return gson{
		"slaveName":    slave.name,
		"visitorCount": server.RoomByName("room.slave.all.user." + slave.name).Count(),
		"state":        slave.state,
		"startTime":    slave.startTime,
		"endTime":      slave.endTime,
		"user": func() interface{} {
			if slave.userWaitState != nil {
				return gson{
					"order":     slave.userWaitState.order,
					"sessionId": slave.userWaitState.session.Id(),
				}
			}
			return nil
		}(),
	}
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
