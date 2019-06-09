package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
)

var slaveRpcMap = map[string]realtime.RpcHandler{
	"slave.hello": realtime.RpcHandleFunc(func(ss *realtime.Session, p gson) (result interface{}, err error) {
		return
	}),
	"slave.session.public.info": sessionPublicInfo,
}

var slaveEventMap = map[string]realtime.EventHandler{
	"slave.broadcast": roleBroadcast,
	"slave.state.change": realtime.EventHandleFunc(func(ss *realtime.Session, data interface{}) {
		slave := ss.Data.(*roleSlave)
		state := data.(gson)["state"].(string)
		if state == "running" {
			slave.state = slaveStateRunning
			slave.startTime = int(data.(gson)["startTime"].(float64))
			slave.endTime = int(data.(gson)["endTime"].(float64))
		}
		if state == "shutting" {
			slave.state = slaveStateWait
		}
	}),
}
