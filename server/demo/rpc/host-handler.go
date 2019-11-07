package rpc

import (
	"errors"
	"strings"

	"github.com/peterq/pan-light/server/realtime"
)

var hostRpcMap = map[string]realtime.RpcHandler{
	"host.rtc.candidate": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		userSessionId := p["sessionId"].(string)
		requestId := p["requestId"].(string)
		user, ok := server.SessionById(realtime.SessionId(userSessionId))
		if !ok {
			err = errors.New("user not in here")
		}
		user.Emit("host.candidate.ok", gson{
			"candidate": p["candidate"],
			"sessionId": userSessionId,
			"requestId": requestId,
		})
		return
	}),
	"host.next.user": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		host := ss.Data.(*roleHost)
		slaveName := p["slave"].(string)
		if strings.Index(slaveName, host.name) != 0 {
			err = errors.New("forbidden")
		}
		slave := host.slaves[slaveName]
		manager.waitSessionMapLock.Lock()
		defer manager.waitSessionMapLock.Unlock()
		for i := manager.lastInServiceOrder + 1; i <= manager.lastDistributedOrder; i++ {
			state, ok := manager.waitSessionMap[i]
			if ok {
				delete(manager.waitSessionMap, i)
				result = gson{
					"order":     state.order,
					"ticket":    state.ticket,
					"sessionId": state.session.Id(),
				}
				manager.lastInServiceOrder = state.order
				server.RoomByName("room.slave.all.user." + slaveName).Join(state.session.Id())
				server.RoomByName("room.all.user").Broadcast("ticket.turn", gson{
					"order":     state.order,
					"sessionId": state.session.Id(),
					"host":      host.name,
					"slave":     slaveName,
				})
				slave.userWaitState = state
				slave.state = slaveStateStarting
				return
			}
		}
		err = errors.New("无用户在排队")
		return
	}),
	"host.slave.register": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		host := ss.Data.(*roleHost)
		if host.slaves != nil {
			err = errors.New("已经注册过了")
			return
		}
		host.slaves = map[string]*roleSlave{}
		for _, name := range p["slaves"].([]interface{}) {
			slave := &roleSlave{
				name:  name.(string),
				host:  host,
				state: slaveStateWait,
			}
			host.slaves[name.(string)] = slave
			manager.slaveMap[name.(string)] = slave
		}
		host.wsAgentUrl = p["ws_agent_url"].(string)
		return
	}),
	"host.hello": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		return
	}),
}

var hostEventMap = map[string]realtime.EventHandler{
	"host": realtime.EventHandleFunc(func(ss *realtime.Session, data interface{}) {

	}),
	"host.slave.exit": realtime.EventHandleFunc(func(ss *realtime.Session, data interface{}) {
		host := ss.Data.(*roleHost)
		slaveName := data.(string)
		slave := host.slaves[slaveName]

		unexpected := slave.state == slaveStateRunning
		slave.state = slaveStateWait
		room := server.RoomByName("room.slave.all.user." + slaveName)
		room.Broadcast("slave.exit", gson{
			"unexpected": unexpected,
		})
		for _, member := range room.Members() {
			room.Remove(member)
		}
	}),
	"host.broadcast": roleBroadcast,
}
