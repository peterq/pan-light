package host

import (
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

func startWsAgentServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello ws agent"))
	})
	http.HandleFunc("/ws", beforeWsAgent)
	log.Println("ws agent server port", host.wsAgentPort)

	_, e1 := ioutil.ReadFile("cert.pem")
	_, e2 := ioutil.ReadFile("key.pem")
	useHttps := e1 == nil && e2 == nil
	log.Println("ws agent server port", host.wsAgentPort, useHttps)
	if useHttps {
		http.ListenAndServeTLS(":"+host.wsAgentPort, "cert.pem", "key.pem", nil)
	} else {
		http.ListenAndServe(":"+host.wsAgentPort, nil)
	}
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
