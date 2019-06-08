package pan_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/peterq/pan-light/pc/storage"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"
)

var urlMap = map[string]string{
	"home":  "https://pan.baidu.com/disk/home",
	"list":  "https://pan.baidu.com/api/list",
	"dlink": "https://pan.baidu.com/api/download",
	"usage": "https://pan.baidu.com/api/quota",
}

func newRequest(method, url string) *http.Request {
	if u, ok := urlMap[url]; ok {
		url = u
	}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("user-agent", BaiduUA)
	if err != nil {
		log.Println(err)
	}
	return req
}

func GetSign() (ctx map[string]interface{}, err error) {
	if storage.UserState.PanCookie != nil {
		cookieJar, _ = cookiejar.New(nil)
		httpClient.Jar = cookieJar

		var cookies []*http.Cookie
		for _, c := range storage.UserState.PanCookie {
			cookies = append(cookies, &http.Cookie{
				Name:   c.Key,
				Value:  c.Value,
				Domain: ".baidu.com",
			})
		}
		u, _ := url.Parse("https://pan.baidu.com")
		cookieJar.SetCookies(u, cookies)
	}
	req := newRequest("GET", "home")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(req.Cookies())
		log.Println(err)
		return
	}
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(string(bin))
		return
	}
	body := string(bin)
	reg := regexp.MustCompile(`var context=(.*);\n`)
	find := reg.FindStringSubmatch(body)
	if res.Request.URL.String() != urlMap["home"] {
		log.Println("重定向: ", res.Request.URL.String())
		err = errors.New("未登录")
		return
	}
	raw := tJson{}
	err = json.Unmarshal([]byte(find[1]), &raw)
	if err != nil {
		log.Println(body)
		return
	}
	handleLoginSession(&raw)
	ctx = raw
	return
}

// 获取网盘使用空间
func Usage() (result interface{}, err error) {
	req := newRequest("GET", "usage")
	params := map[string]interface{}{
		"checkexpire": 1,
		"checkfree":   1,
		"channel":     "chunlei",
		"web":         1,
		"app_id":      250528,
		"bdstoken":    LoginSession.Bdstoken,
		"logid":       time.Now().UnixNano(),
		"clienttype":  0,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取磁盘空间错误, 错误码" + fmt.Sprint(data["errno"]))
	}
	result = data
	return
}

// 获取目录下的文件(夹)
func ListDir(path string) (list interface{}, err error) {
	req := newRequest("GET", "list")
	params := map[string]interface{}{
		"channel":    "chunlei",
		"clienttype": 0,
		"web":        1,
		"showempty":  1,
		"num":        10000,
		"t":          time.Now().Unix() * 1000,
		"dir":        path,
		"page":       1,
		"desc":       1,
		"order":      "name",
		"_":          time.Now().Unix() * 1000,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取文件夹信息错误, 错误码" + fmt.Sprint(data["errno"]))
	}
	list = data["list"]
	return
}

// 链接解析
func Link(fid string) (link string, err error) {

	if c, ok := linkCacheMap[fid]; !ok {
		linkCacheMap[fid] = fidLinks{}
	} else {
		if c.direct != nil && !c.direct.expired() {
			return c.direct.link, nil
		}
	}

	req := newRequest("GET", "dlink")
	params := map[string]interface{}{
		"sign":       LoginSession.Sign,
		"timestamp":  LoginSession.Timestamp,
		"fidlist":    "[" + fid + "]",
		"type":       "dlink",
		"channel":    "chunlei",
		"web":        1,
		"app_id":     "250528",
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	log.Println(req.URL.String())
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取文件信息错误, 错误码" + fmt.Sprint(data["errno"]))
		return
	}
	link = data["dlink"].([]interface{})[0].(map[string]interface{})["dlink"].(string)
	link = getRedirectedLink(link)
	linkCacheMap[fid] = fidLinks{
		direct: &linkTime{
			link: link,
			time: time.Now(),
		},
	}
	return
}

// vip 转存解析
func linkByVip() (link string, err error) {
	return
}

func getRedirectedLink(link string) string {
	req := newRequest("GET", link)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	end := resp.Request.URL.String()
	log.Println(end)
	resp.Body.Close()
	return end
}
