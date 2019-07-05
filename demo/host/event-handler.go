package host

import "log"

type gson = map[string]interface{}

var eventHandlers = map[string]func(data interface{}, room string){
	"user.connect.request": func(data interface{}, room string) {
		p := data.(gson)
		candidate := p["candidate"].(string)
		requestId := p["requestId"].(string)
		sessionId := p["sessionId"].(string)
		handleNewUser(candidate, sessionId, requestId)
	},
	"session.new": func(data interface{}, room string) {
		startServe()
	},
	"wait.user.new": func(data interface{}, room string) {
		for _, holder := range host.holderMap {
			holder.CheckUser()
		}
	},
	"broadcast.slave": func(data interface{}, room string) {
		p := data.(gson)

		// room for host <-> slaves
		if room == "room.host.slaves."+host.name {
			payload := p["payload"].(gson)
			order := int64(payload["order"].(float64))
			slaveName := payload["slave"].(string)
			holder, ok := host.holderMap[slaveName]
			if !ok {
				log.Println("slave not exist:", slaveName)
				return
			}
			if order != holder.Order() {
				log.Println("order not match:", order, holder.Order())
			}
			log.Printf("%#v", p)
			holder.HandleEvent(p["event"].(string), payload)
		}
	},
}
