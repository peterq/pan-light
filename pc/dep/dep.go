package dep

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/debug"
)

var Fatal = func(str string) {
	debug.PrintStack()
	log.Fatal(str)
}

var initCb []func()

func OnInit(cb func()) {
	initCb = append(initCb, cb)
}
func DoInit() {
	for _, cb := range initCb {
		cb()
	}
	initCb = nil // 防止二次调用
}

var closeCb []func()

func OnClose(cb func()) {
	closeCb = append(closeCb, cb)
}

func DoClose() {
	for _, cb := range closeCb {
		cb()
	}
	closeCb = nil
}

func Reboot() {
	if runtime.GOOS == "windows" {
		ioutil.WriteFile(DataPath("reboot"), []byte("true"), 0664)
	}
	DoClose()
	os.Exit(2)
}

var NotifyQml = func(event string, data map[string]interface{}) {
	log.Println("not ready")
}
