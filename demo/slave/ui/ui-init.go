package ui

import (
	"log"
	"path/filepath"
	"plugin"
)

var endTime int
var nickname string

type gson = map[string]interface{}

func Init(t int, name string) {
	endTime = t
	nickname = name
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	p, err := plugin.Open("./lib/gui-plugin.so")
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

	RegisterAsync(AsyncRouteRegister.(func(routes map[string]func(map[string]interface{},
		func(interface{}), func(interface{}), func(interface{}), chan interface{}))))

	RegisterSync(SyncRouteRegister.(func(routes map[string]func(map[string]interface{}) interface{})))

	StartGui.(func(rccFile, mainQml string))("./lib/qml.rcc", "qrc:/demo/main.qml")
}

var shutdownMagChan chan string

func Shutdown(msg string) {
	shutdownMagChan <- msg
}

func init() {
	shutdownMagChan = make(chan string)
	syncMap(map[string]syncHandler{
		"conf": func(p map[string]interface{}) interface{} {
			return gson{
				"endTime":  endTime,
				"nickname": nickname,
			}
		},
	})
	asyncMap(map[string]asyncHandler{
		"shutdownMsg": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
			resolve(<-shutdownMagChan)
		},
	})
}

func abs(t string) string {
	p, _ := filepath.Abs("./deploy/linux/" + t)
	log.Println(p)
	return p
}
