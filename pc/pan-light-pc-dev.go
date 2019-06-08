// +build plugin

package main

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/functions"
	"log"
	"plugin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
