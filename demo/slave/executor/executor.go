package executor

import (
	"github.com/peterq/pan-light/demo/realtime"
	"github.com/peterq/pan-light/demo/slave/ui"
	"github.com/peterq/pan-light/demo/slave/vnc-password"
	"log"
	"os"
	"time"
)

type gson = map[string]interface{}

type executor struct {
	hostName      string
	slaveName     string
	userSessionId string
	order         int64
	rtOkCh        chan bool
}

func (e *executor) startX() {
	defer rt.Emit("state.change", gson{
		"state": "shutting",
	})
	log.Println("set password")
	vnc_password.SetPassword(env("vnc_operate_pwd"), env("vnc_view_pwd"))
	<-e.rtOkCh
	e.notifyHost("start.ok", gson{})
	startTime := time.Now()
	endTime := startTime.Add(5 * time.Minute)
	endCh := time.After(5 * time.Minute)
	rt.Emit("state.change", gson{
		"state":     "running",
		"startTime": startTime.Unix(),
		"endTime":   endTime.Unix(),
	})
	rt.RegisterEventListener(map[string]func(data interface{}, room string){
		"room.member.remove": func(data interface{}, room string) {
			if room == "room.slave.all.user."+e.slaveName {
				if data.(string) == e.userSessionId {
					log.Println("operator.leave")
					rt.Broadcast(room, "operator.leave", nil)
					e.shutdown("该demo的操作用户已离开, 即将关闭本demo")
				}
			}
		},
	})
	infoMap, err := rt.Call("session.public.info", gson{
		"sessionIds": []string{e.userSessionId},
	})
	nickname := "nickname"
	if err != nil {
		nickname = "nickname err: " + err.Error()
	} else {
		nickname = infoMap.(gson)[e.userSessionId].(gson)["nickname"].(string)
	}
	// 启动面板
	go func() {
		ui.Init(int(endTime.Unix()), nickname)
	}()
	log.Println(os.Environ())
	// 启动程序
	go func() {
		for {
			runCmd("./deploy/linux", "sh", "./pc.sh")
			time.Sleep(time.Second)
			if endTime.Sub(time.Now()) < 10 {
				break
			}
		}
	}()
	// 等待体验结束
	<-endCh
	e.shutdown("体验结束, 即将关闭本demo")
}

func (e *executor) notifyHost(event string, data gson) {
	data["order"] = e.order
	data["slave"] = e.slaveName
	rt.Broadcast("room.host.slaves."+e.hostName, event, data)
}

func (e *executor) shutdown(msg string) {
	ui.Shutdown(msg)
	time.Sleep(5 * time.Second)
	rt.Emit("state.change", gson{
		"state": "shutting",
	})
	os.Exit(0)
}

var rt *realtime.RealTime
var exe *executor
