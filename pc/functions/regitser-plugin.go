// +build plugin

package functions

import (
	"github.com/peterq/pan-light/pc/dep"
	"log"
)

func init() {
	dep.NotifyQml = NotifyQml
}

type syncHandler func(p map[string]interface{}) interface{}

type asyncHandler func(p map[string]interface{}, resolve func(interface{}),
	reject func(interface{}), progress func(interface{}), qmlMsg chan interface{})

var syncRoutes []map[string]syncHandler
var asyncRoutes []map[string]asyncHandler

func syncMap(r map[string]syncHandler) {
	syncRoutes = append(syncRoutes, r)
}

func asyncMap(r map[string]asyncHandler) {
	asyncRoutes = append(asyncRoutes, r)
}

func RegisterAsync(regitser func(routes map[string]func(map[string]interface{},
	func(interface{}), func(interface{}), func(interface{}), chan interface{}))) {

	for _, r := range asyncRoutes {

		r1 := map[string]func(map[string]interface{},
			func(interface{}), func(interface{}), func(interface{}), chan interface{}){}
		for p, h := range r {
			r1[p] = h
		}
		regitser(r1)
	}
}

func RegisterSync(register func(routes map[string]func(map[string]interface{}) interface{})) {

	for _, r := range syncRoutes {

		r1 := map[string]func(map[string]interface{}) interface{}{}
		for p, h := range r {
			r1[p] = h
		}
		register(r1)
	}
}

var NotifyQml = func(event string, data map[string]interface{}) {
	log.Println("this function should be load from plugin", event, data)
}
