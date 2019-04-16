package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
	"github.com/pkg/errors"
)

var userRpcMap = map[string]realtime.RpcHandler{
	"user.hosts.info": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
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
	"user.hosts.hello": realtime.RpcHandleFunc(func(ss *realtime.Session, data gson) (result interface{}, err error) {
		return
	}),
}

var userEventMap = map[string]realtime.EventHandler{
	"user.chat.msg": realtime.EventHandleFunc(func(ss *realtime.Session, data interface{}) {
		payload := data.(gson)
		room := payload["room"].(string)
		msg := payload["msg"]
		if ss.InRoom(room) {
			server.RoomByName(room).Broadcast("chat.msg.new", gson{
				"from": ss.Id(),
				"msg":  msg,
				"room": room,
			}, ss.Id())
		}
	}),
}
