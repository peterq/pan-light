package pan_viper

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type gson = map[string]interface{}

const baiduUa = "netdisk;2.2.3;pc;pc-mac;10.14.5;macbaiduyunguanjia"

func makeHttpClient(cookieStr string) (http.Client, string) {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://pan.baidu.com")
	var cookies []*http.Cookie
	parts := strings.Split(strings.TrimSpace(cookieStr), ";")
	var bduss string
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
		if len(parts[i]) == 0 {
			continue
		}
		name, val := parts[i], ""
		if j := strings.Index(name, "="); j >= 0 {
			name, val = name[:j], name[j+1:]
		}
		if name == "BDUSS" {
			bduss = val
		}
		cookies = append(cookies, &http.Cookie{Name: name, Value: val, Domain: ".baidu.com"})
	}
	jar.SetCookies(u, cookies)
	httpClient := http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Referer")
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
		Jar:     jar,
		Timeout: 15 * time.Second,
	}
	return httpClient, bduss
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
