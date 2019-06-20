// +build !plugin

package functions

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/gui/bridge"
)

func init() {
	dep.NotifyQml = NotifyQml
}

type syncHandler func(p map[string]interface{}) (result interface{})

type asyncHandler func(p map[string]interface{},
	resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{})

func syncMap(r map[string]syncHandler) {
	r1 := map[string]func(map[string]interface{}) interface{}{}
	for p, h := range r {
		r1[p] = h
	}
	bridge.SyncRouteRegister(r1)
}

func asyncMap(r map[string]asyncHandler) {
	r1 := map[string]func(map[string]interface{},
		func(interface{}), func(interface{}), func(interface{}), chan interface{}){}
	for p, h := range r {
		r1[p] = h
	}
	bridge.AsyncRouteRegister(r1)
}

func NotifyQml(event string, data map[string]interface{}) {
	bridge.NotifyQml(event, data)
}
