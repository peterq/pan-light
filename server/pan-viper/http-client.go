package pan_viper

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type gson = map[string]interface{}

const baiduUa = "netdisk;4.6.2.0;PC;PC-Windows;10.0.10240;WindowsBaiduYunGuanJia"

func makeHttpClient(bduss string) http.Client {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://pan.baidu.com")
	jar.SetCookies(u, []*http.Cookie{
		{
			Name:     "BDUSS",
			Value:    bduss,
			Path:     "/",
			HttpOnly: false,
			MaxAge:   math.MaxInt32,
			Domain:   "*.baidu.com",
		},
	})
	httpClient := http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Referer")
			log.Println(req.URL)
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
		Jar:     jar,
		Timeout: 15 * time.Second,
	}
	return httpClient
}

func newRequest(method, link string, body ...gson) *http.Request {
	var bd io.Reader
	if len(body) == 1 {
		formData := url.Values{}
		for key, value := range body[0] {
			formData.Add(key, fmt.Sprint(value))
		}
		bd = strings.NewReader(formData.Encode())
	}
	req, err := http.NewRequest(method, link, bd)
	req.Header.Set("user-agent", baiduUa)
	req.Header.Set("referer", "https://pan.baidu.com")
	if err != nil {
		log.Println(err)
	}
	return req
}
