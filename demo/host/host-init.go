package host

import (
	"fmt"
	"github.com/peterq/pan-light/demo/realtime"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var host = &struct {
	name       string
	password   string
	wsAddr     string
	initLock   sync.Mutex
	inited     bool
	slaveCount int
	slaves     []string
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
		return
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
	log.Println(rt.Call("host.next.user", gson{
		"slave": host.slaves[0],
	}))
}
