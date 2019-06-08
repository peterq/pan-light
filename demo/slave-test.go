package main

import (
	"fmt"
	"github.com/peterq/pan-light/demo/slave/ui"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"
)

func main() {
	pidFile := "slave-test.pid"
	bin, e := ioutil.ReadFile(pidFile)
	if e == nil {
		pid, _ := strconv.Atoi(string(bin))
		syscall.Kill(pid, syscall.SIGKILL)
	}
	ioutil.WriteFile(pidFile, []byte(fmt.Sprint(os.Getpid())), os.ModePerm)

	go func() {
		time.Sleep(5 * time.Second)
		ui.Shutdown("操作该demo的用户已离开, 即将关闭本demo...")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
	go func() {
		ui.Init(int(time.Now().Add(5*time.Minute).Unix()), "风清扬")
	}()
	select {}
}
