package host

import "log"

type gson = map[string]interface{}

var eventHandlers = map[string]func(data interface{}){
	"user.connect.request": func(data interface{}) {
		p := data.(gson)
		candidate := p["candidate"]
		requestId := p["requestId"].(string)
		sessionId := p["sessionId"].(string)
		log.Println(candidate, requestId, sessionId)
	},
	"session.new": func(data interface{}) {
		startServe()
	},
}
