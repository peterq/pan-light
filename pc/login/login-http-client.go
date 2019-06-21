package login

import (
	"encoding/json"
	"github.com/peterq/pan-light/pc/storage"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

var cookieJar *cookiejar.Jar
var httpClient http.Client

var urlMap = map[string]string{
	"gotoWx":   "https://passport.baidu.com/phoenix/account/startlogin?type=42&tpl=netdisk&u=https%3A%2F%2Fpan.baidu.com%2Fdisk%2Fhome&display=page&act=implicit&subpro=netdisk_web",
	"gotoQQ":   "https://passport.baidu.com/phoenix/account/startlogin?type=15&tpl=netdisk&u=https%3A%2F%2Fpan.baidu.com%2Fdisk%2Fhome&display=page&act=implicit&subpro=netdisk_web",
	"gotoSina": "https://passport.baidu.com/phoenix/account/startlogin?type=2&tpl=netdisk&u=https%3A%2F%2Fpan.baidu.com%2Fdisk%2Fhome&display=page&act=implicit&subpro=netdisk_web",
}

func newRequest(method, url string) *http.Request {
	if u, ok := urlMap[url]; ok {
		url = u
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println(err)
	}
	return req
}

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
		Timeout: 60 * time.Second,
	}
}

func readHtml(reader io.Reader) string {
	html, _ := ioutil.ReadAll(reader)
	return string(html)
}

type tJson map[string]interface{}
type tBin []byte

func handleLoginSuccess() (err error) {
	link := "https://pan.baidu.com/disk/home"
	u, _ := url.Parse(link)
	log.Println("登录成功", cookieJar.Cookies(u))
	//resp, _ := httpClient.Get(link)
	//log.Println(readHtml(resp.Body))
	req := newRequest("GET", link)
	res, err := httpClient.Do(req)
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(string(bin))
		return
	}
	body := string(bin)
	reg := regexp.MustCompile(`var context=(.*);\n`)
	find := reg.FindStringSubmatch(body)
	if len(find) != 2 {
		return errors.New("未找到context")
	}
	raw := tJson{}
	err = json.Unmarshal([]byte(find[1]), &raw)
	if err != nil {
		log.Println(body)
		return
	}
	log.Println(raw["username"])
	storage.OnLogin(raw["username"].(string))
	log.Println(req.Cookies())
	var cookies []*storage.Cookies
	for _, c := range cookieJar.Cookies(u) {
		cookies = append(cookies, &storage.Cookies{
			Key:   c.Name,
			Value: c.Value,
		})
	}
	storage.UserState.PanCookie = cookies
	log.Println("global pan cookie", storage.UserState.PanCookie)
	return
}
