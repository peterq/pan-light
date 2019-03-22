// +build plugin

package main

import (
	"github.com/peterq/pan-light/pc/functions"
	"plugin"
)

func main() {
	p, err := plugin.Open("./gui/gui-plugin.so")
	if err != nil {
		panic(err)
	}
	StartGui, err := p.Lookup("StartGui")
	if err != nil {
		panic(err)
	}
	SyncRouteRegitser, err := p.Lookup("SyncRouteRegitser")
	if err != nil {
		panic(err)
	}
	AsyncRouteRegitser, err := p.Lookup("AsyncRouteRegitser")
	if err != nil {
		panic(err)
	}

	functions.RegisterAsync(AsyncRouteRegitser.(func(routes map[string]func(map[string]interface{},
		func(interface{}), func(interface{}), func(interface{}), chan interface{}))))

	functions.RegisterSync(SyncRouteRegitser.(func(routes map[string]func(map[string]interface{}) interface{})))

	StartGui.(func(rccFile, mainQml string))("./gui/qml/qml.rcc", "qrc:/main.qml")
}
