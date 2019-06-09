package executor

import (
	"github.com/peterq/pan-light/demo/realtime"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func Start() {
	log.Println("hello pan light, real_time connecting")
	pwd := env("host_password")
	os.Unsetenv("host_password")
	rt = &realtime.RealTime{
		WsAddr:       env("ws_addr"),
		Role:         "slave",
		HostName:     env("host_name"),
		HostPassWord: pwd,
		SlaveName:    env("slave_name"),
		OnConnected:  nil,
	}
	rt.Init()
	order, _ := strconv.ParseInt(env("demo_order"), 10, 64)
	exe = &executor{
		hostName:      rt.HostName,
		slaveName:     rt.SlaveName,
		order:         order,
		userSessionId: env("demo_user"),
	}
	exe.startX()
}

func cmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func runCmd(path, name string, arg ...string) {
	c := cmd(name, arg...)
	c.Dir, _ = filepath.Abs(path)
	c.Run()
}

func env(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic("this env var must be set: " + key)
	}
	return v
}
