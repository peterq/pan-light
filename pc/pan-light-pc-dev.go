// +build plugin

package main

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/functions"
	"log"
	"os"
	"os/exec"
	"plugin"
	"syscall"
)

const startCmd = "pan_light_start"

func master() {
	log.Println("master process")
START_PAN:
	c := exec.Command(os.Args[0], os.Args[1:]...)
	c.Args[0] = startCmd
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	err := c.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code := status.ExitStatus()
				if code == 2 {
					goto START_PAN
				}
			}
		}
	}
	log.Fatal(err)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if os.Args[0] != startCmd {
		master()
	}
	log.Println("pan-light process")
	defer func() {
		dep.DoClose()
	}()
	dep.DoInit()
	p, err := plugin.Open("./gui/gui-plugin.so")
	if err != nil {
		panic(err)
	}
	StartGui, err := p.Lookup("StartGui")
	if err != nil {
		panic(err)
	}
	SyncRouteRegister, err := p.Lookup("SyncRouteRegister")
	if err != nil {
		panic(err)
	}
	AsyncRouteRegister, err := p.Lookup("AsyncRouteRegister")
	if err != nil {
		panic(err)
	}
	NotifyQml, err := p.Lookup("NotifyQml")
	if err != nil {
		panic(err)
	}

	functions.RegisterAsync(AsyncRouteRegister.(func(routes map[string]func(map[string]interface{},
		func(interface{}), func(interface{}), func(interface{}), chan interface{}))))

	functions.RegisterSync(SyncRouteRegister.(func(routes map[string]func(map[string]interface{}) interface{})))

	functions.NotifyQml = NotifyQml.(func(event string, data map[string]interface{}))
	dep.NotifyQml = functions.NotifyQml
	StartGui.(func(rccFile, mainQml string))("./gui/qml/qml.rcc", "qrc:/main.qml")
}
