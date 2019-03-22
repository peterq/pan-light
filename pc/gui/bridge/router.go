package bridge

import (
	"context"
	"fmt"
	"github.com/peterq/pan-light/pc/gui/qt-rpc"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

// 业务逻辑侧衔接ui

var sendDataToQml func(string)
var OnClose func(func())
var UiContext = context.WithValue(context.Background(), "start_time", time.Now())

type tJson map[string]interface{}

func NotifyQml(event string, data map[string]interface{}) {
	qt_rpc.NotifyQml(event, (*qt_rpc.Gson)(&data))
}

var AsyncTaskChanMap = make(map[string]chan interface{})
var AsyncTaskChanMapLock = new(sync.RWMutex)

func init() {
	qt_rpc.CallGoAsync = callGoAsync
	qt_rpc.CallGoSync = callGoSync
}

func callGoAsync(data *qt_rpc.Gson) {
	var err error
	p := *data
	action := p["action"].(string)
	if fn, ok := callAsyncMap[action]; ok {
		go func() {
			finish := func() {}
			defer func() {
				if e := recover(); e != nil {
					err = e.(error)
					log.Println("调用go api函数出错"+action, e)
					log.Printf("stack %s", debug.Stack())
					NotifyQml("call.ret", map[string]interface{}{
						"type":   "reject",
						"callId": p["callId"],
						"reject": err.Error(),
					})
					finish()
				}
			}()
			var qmlMsg chan interface{}
			if withCh, ok := p["chan"]; ok && withCh.(bool) {
				qmlMsg = make(chan interface{})
				AsyncTaskChanMapLock.Lock()
				AsyncTaskChanMap[p["callId"].(string)] = qmlMsg
				AsyncTaskChanMapLock.Unlock()
				finish = func() {
					AsyncTaskChanMapLock.Lock()
					defer AsyncTaskChanMapLock.Unlock()
					close(qmlMsg)
					delete(AsyncTaskChanMap, p["callId"].(string))
				}
			}
			fn(p["param"].(map[string]interface{}), func(i interface{}) {
				NotifyQml("call.ret", map[string]interface{}{
					"type":    "resolve",
					"callId":  p["callId"],
					"resolve": i,
				})
				finish()
			}, func(i interface{}) {
				if e, ok := i.(error); ok {
					i = e.Error()
				}
				NotifyQml("call.ret", map[string]interface{}{
					"type":   "reject",
					"callId": p["callId"],
					"reject": i,
				})
				finish()
			}, func(i interface{}) {
				NotifyQml("call.ret", map[string]interface{}{
					"type":     "progress",
					"callId":   p["callId"],
					"progress": i,
				})
			}, qmlMsg)
		}()
	} else {
		NotifyQml("call.ret", map[string]interface{}{
			"type":   "reject",
			"callId": p["callId"],
			"reject": fmt.Sprintf("api [%s] not exist", action),
		})
	}
}

func callGoSync(data *qt_rpc.Gson) (retPointer *qt_rpc.Gson) {
	ret := qt_rpc.Gson{}
	defer func() { retPointer = &ret }()
	p := *data
	var err error
	action := p["action"].(string)
	if fn, ok := callSyncMap[action]; ok {
		defer func() {
			if e := recover(); e != nil {
				err = e.(error)
				log.Println("调用go api函数出错"+action, e)
				log.Printf("stack %s", debug.Stack())
				ret["error"] = "call go panic"
				return
			}
		}()
		result := fn(p["param"].(map[string]interface{}))
		ret["result"] = result
		return
	} else {
		ret["error"] = "api not exist"
		return
	}
}
