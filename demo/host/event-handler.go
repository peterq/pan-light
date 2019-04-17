package host

type gson = map[string]interface{}

var eventHandlers = map[string]func(data interface{}){
	"user.connect.request": func(data interface{}) {
		p := data.(gson)
		candidate := p["candidate"].(string)
		requestId := p["requestId"].(string)
		sessionId := p["sessionId"].(string)
		handleNewUser(candidate, sessionId, requestId)
	},
	"session.new": func(data interface{}) {
		startServe()
	},
	"wait.user.new": func(data interface{}) {
		for _, holder := range host.holderMap {
			holder.CheckUser()
		}
	},
}
