package pan_api

import (
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/storage"
	"io"
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
	if ca, ok := linkCacheMap[request.URL.Query().Get("fid")]; ok {
		lt := ca.direct

		myReq := newRequest("GET", lt.link)

		for k, vs := range request.Header {
			if k == "Referer" {
				continue
			}
			for _, h := range vs {
				//log.Println(k, h)
				myReq.Header.Add(k, h)
			}
		}
		//log.Println("-----------------")
		myReq.Header.Set("user-agent", BaiduUA)

		resp, err := httpClient.Do(myReq)
		if err != nil {
			log.Println(err)
			return
		}
		for k, vs := range resp.Header {
			if k == "Content-Disposition" {
				continue
			}
			for _, h := range vs {
				//log.Println(k, h)
				writer.Header().Add(k, h)
			}
			writer.Header().Set("Connection", "close")
		}
		writer.WriteHeader(resp.StatusCode)
		io.Copy(writer, resp.Body)
	} else {
		writer.WriteHeader(500)
	}
}
