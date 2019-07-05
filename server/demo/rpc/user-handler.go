package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
	"github.com/pkg/errors"
)

var userRpcMap = map[string]realtime.RpcHandler{
	"user.hosts.info": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		manager.hostMapLock.RLock()
		defer manager.hostMapLock.RUnlock()
		var arr []gson
		for _, host := range manager.hostMap {
			var slaves []gson
			for slaveName := range host.slaves {
				slaves = append(slaves, gson{
					"slaveName":    slaveName,
					"visitorCount": server.RoomByName("room.slave.all.user." + slaveName).Count(),
				})
			}
			arr = append(arr, gson{
				"name":       host.name,
				"wsAgentUrl": host.wsAgentUrl,
				"slaves":     slaves,
			})
		}
		return arr, nil
	}),
	"user.host.detail": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		manager.hostMapLock.RLock()
		defer manager.hostMapLock.RUnlock()
		host, ok := manager.hostMap[p["hostName"].(string)]
		if !ok {
			err = errors.New("host not exist")
			return
		}
		var slaves []gson
		for _, slave := range host.slaves {
			slaves = append(slaves, slave.publicInfo())
		}
		result = gson{
			"slaves": slaves,
		}
		return
	}),
	"user.ping": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		return "pong", nil
	}),
	"user.connect.host": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		candidate := data["candidate"]
		requestId := data["requestId"].(string)
		hostName := data["hostName"].(string)
		manager.hostMapLock.Lock()
		defer manager.hostMapLock.Unlock()
		host, ok := manager.hostMap[hostName]
		if !ok {
			err = errors.New("host 不存在")
			return
		}
		host.session.Emit("user.connect.request", gson{
			"candidate": candidate,
			"requestId": requestId,
			"sessionId": ss.Id(),
		})
		return
	}),
	// 请求在线体验票据
	"user.ticket.new": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		user := ss.Data.(*roleUser)
		return user.requestTicket()
	}),
	"user.hosts.hello": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		return
	}),
	"user.join.slave": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		manager.slaveMapLock.RLock()
		defer manager.slaveMapLock.RUnlock()
		slaveName := data["slave"].(string)
		_, ok := manager.slaveMap[slaveName]
		if !ok {
			err = errors.New("slave 不存在")
		}
		roomName := "room.slave.all.user." + slaveName
		server.RoomByName(roomName).Join(ss.Id())
		return
	}),
	"user.leave.slave": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		manager.slaveMapLock.RLock()
		defer manager.slaveMapLock.RUnlock()
		slaveName := data["slave"].(string)
		_, ok := manager.slaveMap[slaveName]
		if !ok {
			err = errors.New("slave 不存在")
		}
		roomName := "room.slave.all.user." + slaveName
		server.RoomByName(roomName).Remove(ss.Id())
		return
	}),
	"user.room.members": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		roomName := data["room"].(string)
		if server.RoomExist(roomName) {
			return server.RoomByName(roomName).Members(), nil
		}
		return nil, errors.New("room not exist")
	}),
	"user.session.public.info": sessionPublicInfo,
}

var userEventMap = map[string]realtime.EventHandler{
	"user.broadcast": roleBroadcast,
}
