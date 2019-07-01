package host

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

func startWsAgentServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello ws agent"))
	})
	http.HandleFunc("/ws", beforeWsAgent)
	log.Println("ws agent server port", host.wsAgentPort)
	http.ListenAndServe(":"+host.wsAgentPort, nil)
}

func beforeWsAgent(writer http.ResponseWriter, request *http.Request) {
	slaveName := request.URL.Query().Get("slave")
	holder, ok := host.holderMap[slaveName]
	if !ok {
		writer.Write([]byte("slave not exist"))
		return
	}
	websocket.Handler(holder.WsProxy).ServeHTTP(writer, request)
}
