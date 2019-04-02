package pan_api

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
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

var BaiduUA = "netdisk;4.6.2.0;PC;PC-Windows;10.0.10240;WindowsBaiduYunGuanJia"

type tBin []byte
type tJson map[string]interface{}

type linkTime struct {
	link string
	time time.Time
}

func (l *linkTime) expired() bool {
	return false
}

type fidLinks struct {
	direct *linkTime
	vip    *linkTime
}

var linkCacheMap = map[string]fidLinks{}
