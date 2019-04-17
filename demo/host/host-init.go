package host

import (
	"context"
	"fmt"
	"github.com/peterq/pan-light/demo/host/instance"
	"github.com/peterq/pan-light/demo/realtime"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var host = &struct {
	name       string
	password   string
	wsAddr     string
	slaveCount int

	slaves []string

	initLock       sync.Mutex
	inited         bool
	cancelServe    context.CancelFunc
	cancelInsServe []context.CancelFunc
	holderMap      map[string]*instance.Holder // slave name -> holder

	p2pMap     map[string]*p2p
	p2pMapLock sync.Mutex
}{}
var rt *realtime.RealTime

func Start() {
	host.name = env("host_name")
	host.password = env("host_password")
	host.wsAddr = env("ws_addr")
	var err error
	host.slaveCount, err = strconv.Atoi(env("slave_count"))
	if err != nil {
		panic(err)
	}
	rt = &realtime.RealTime{
		WsAddr:       host.wsAddr,
		Role:         "host",
		HostName:     host.name,
		HostPassWord: host.password,
		OnConnected:  nil,
	}
	rt.Init()
	rt.RegisterEventListener(eventHandlers)
	select {}
}

func env(name string) string {
	s, e := os.LookupEnv("pan_light_" + name)
	if !e {
		panic(fmt.Sprintf("env %s must be set", name))
	}
	return strings.Trim(s, " \"")
}

func startServe() {
	host.initLock.Lock()
	defer host.initLock.Unlock()
	if host.inited {
		host.cancelServe()
		host.inited = false
		host.cancelServe = nil
		host.cancelInsServe = nil
	}
	host.slaves = make([]string, host.slaveCount)
	for i := 0; i < host.slaveCount; i++ {
		host.slaves[i] = host.name + ".slave." + strconv.Itoa(i)
	}
	_, err := rt.Call("host.slave.register", gson{
		"slaves": host.slaves,
	})
	if err != nil {
		log.Println("注册slave失败", err)
		return
	}
	serveCtx, cancel := context.WithCancel(context.Background())
	host.cancelServe = cancel
	host.cancelInsServe = make([]context.CancelFunc, host.slaveCount)
	host.holderMap = map[string]*instance.Holder{}
	host.p2pMap = map[string]*p2p{}
	host.inited = true
	for idx, slaveName := range host.slaves {
		holder := &instance.Holder{
			SlaveName:    slaveName,
			HostName:     host.name,
			HostPassword: host.password,
		}
		host.holderMap[slaveName] = holder
		ctx, cancel := context.WithCancel(serveCtx)
		host.cancelInsServe[idx] = cancel
		time.Sleep(100 * time.Millisecond)
		go holder.Init(rt, ctx)
	}
}
