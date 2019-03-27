package bridge

import (
	"time"
)

var callSyncMap = map[string]func(map[string]interface{}) interface{}{
	// 通知异步任务
	"asyncTaskMsg": func(p map[string]interface{}) interface{} {
		AsyncTaskChanMapLock.RLock()
		ch, ok := AsyncTaskChanMap[p["asyncCallId"].(string)]
		AsyncTaskChanMapLock.RUnlock()
		if !ok {
			return false
		}
		ch <- p["msg"]
		return true
	},
	// 获取当前时间
	"time": func(p map[string]interface{}) interface{} {
		return time.Now().UnixNano()
	},
}

var callAsyncMap = map[string]func(map[string]interface{},
	func(interface{}), func(interface{}), func(interface{}), chan interface{}){

	"wait": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		for i := 0; i < int(p["time"].(float64)); i++ {
			time.Sleep(time.Second)
			progress("current time is " + time.Now().String())
		}
		resolve("wait complete")
	},
}

func SyncRouteRegister(routes map[string]func(map[string]interface{}) interface{}) {
	for path, handler := range routes {
		callSyncMap[path] = handler
	}
}

func AsyncRouteRegister(routes map[string]func(map[string]interface{},
	func(interface{}), func(interface{}), func(interface{}), chan interface{})) {
	for path, handler := range routes {
		callAsyncMap[path] = handler
	}
}
