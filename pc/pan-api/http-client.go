package pan_api

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var cookieJar *cookiejar.Jar
var httpClient http.Client

func init() {
	var e error
	cookieJar, e = cookiejar.New(nil)
	if e != nil {
		panic(e)
	}
	httpClient = http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Referer")
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
		Jar:     cookieJar,
		Timeout: 10 * time.Second,
	}
}

func readHtml(reader io.Reader) string {
	html, _ := ioutil.ReadAll(reader)
	return string(html)
}

func VideoProxy(writer http.ResponseWriter, request *http.Request, targetLink string) {
	myReq := newRequest("GET", targetLink)

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
}

var BaiduUA = "netdisk;2.2.3;pc;pc-mac;10.14.5;macbaiduyunguanjia"

type tBin []byte
type tJson map[string]interface{}
