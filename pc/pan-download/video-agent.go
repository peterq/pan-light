package pan_download

import (
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/pan-api"
	"github.com/peterq/pan-light/pc/storage"
	"log"
	"net"
	"net/http"
	"time"
)

func init() {
	dep.OnInit(func() {
		go startAgentServer()
	})
}

type headChecker struct {
	real http.Handler
}

func (c headChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.real.ServeHTTP(w, r)
}

func startAgentServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/videoAgent", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		videoAgent(writer, request)
	})

	mux.HandleFunc("/exit", func(writer http.ResponseWriter, request *http.Request) {
		dep.Fatal("exit by api")
	})

	if dep.Env.Dev && storage.Global.InternalServerPort > 0 {
		_, err := http.Get("http://127.0.0.1:" + fmt.Sprint(storage.Global.InternalServerPort) + "/exit")
		if err != nil {
			time.Sleep(time.Second)
		}
	}

	log.Println("Listening...")
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		dep.Fatal(err.Error())
	}
	storage.Global.InternalServerPort = int64(listener.Addr().(*net.TCPAddr).Port)
	log.Println("Using port:", storage.Global.InternalServerPort)
	dep.Env.InternalServerUrl = "http://127.0.0.1:" + fmt.Sprint(storage.Global.InternalServerPort)
	err = http.Serve(listener, headChecker{real: mux})
	if err != nil {
		dep.Fatal(err.Error())
	}
}

func videoAgent(writer http.ResponseWriter, request *http.Request) {
	link, err := LinkResolver(request.URL.Query().Get("fid"))
	if err != nil {
		writer.WriteHeader(500)
		log.Println(err)
		return
	}
	pan_api.VideoProxy(writer, request, link)
}
