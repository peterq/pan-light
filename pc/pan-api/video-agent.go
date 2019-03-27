package pan_api

import (
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"io"
	"log"
	"net/http"
	"strings"
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
	ua := r.Header.Get("User-Agent")
	if strings.Index(ua, dep.Env.ElectronSecretUA) < 0 && false {
		w.WriteHeader(404)
		w.Write([]byte("hello, pan-light " + dep.Env.VersionString))
		return
	}
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

	if dep.Env.Dev {
		_, err := http.Get("http://127.0.0.1:" + fmt.Sprint(dep.Env.ListenPort) + "/exit")
		if err != nil {
			time.Sleep(time.Second)
		}
	}

	log.Println("Listening...")
	dep.Env.InternalServerUrl = "http://127.0.0.1:" + fmt.Sprint(dep.Env.ListenPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", dep.Env.ListenPort), headChecker{real: mux})
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
