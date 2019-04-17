package rpc

import (
	"github.com/kataras/iris/core/errors"
	"github.com/peterq/pan-light/server/realtime"
	"strings"
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
	// 作废, 不走这个逻辑
	/*"host.check.ticket": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		order := p["order"].(int64)
		ticket := p["ticket"].(string)
		manager.waitSessionMapLock.RLock()
		defer manager.waitSessionMapLock.RUnlock()
		wait, ok := manager.waitSessionMap[order]
		if wait.ticket != ticket {
			err = errors.New("ticket 错误")
		}
		if !ok {
			err = errors.New("ticket 已失效")
			return
		}
		if !wait.inService {
			err = errors.New("当前用户还未到")
			return
		}
		if wait.serviced {
			err = errors.New("票据已被使用过")
		}
		return
	}),*/
	"host.next.user": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		host := ss.Data.(*roleHost)
		slave := p["slave"].(string)
		if strings.Index(slave, host.name) != 0 {
			err = errors.New("forbidden")
		}
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
				name: name.(string),
				host: host,
			}
			host.slaves[name.(string)] = slave
		}
		return
	}),
	"host.hello": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		return
	}),
}

var hostEventMap = map[string]realtime.EventHandler{
	"host": realtime.EventHandleFunc(func(ss *realtime.Session, data interface{}) {

	}),
}
