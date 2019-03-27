package qt_rpc

import (
	"log"
	"time"
)

type Gson map[string]interface{}

var CallGoSync = func(gson *Gson) *Gson {
	return &Gson{
		"error": "go rpc service not initialized",
	}
}

var CallGoAsync = func(gson *Gson) {
	go func() {
		result := &Gson{
			"type":   "reject",
			"callId": (*gson)["callId"],
			"reject": "go rpc service not initialized",
		}
		NotifyQml("call.ret", result)
	}()
}

var NotifyQml = func(event string, data *Gson) {
	log.Println("qml msg handler not initialized", event, *data)
}

func init() {
	go func() {
		time.Sleep(1 * time.Second)
		NotifyQml("fuck", &Gson{"name": "trump"})
	}()
}
