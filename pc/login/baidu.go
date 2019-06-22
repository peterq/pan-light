package login

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type BaiduLoginOption struct {
	Ctx       context.Context
	OnError   func(err error)
	OnQrCode  func(img string, pageUrl string)
	OnScan    func()
	OnConfirm func()
	OnSuccess func()
}

func BaiduLogin(option *BaiduLoginOption) {
	var err error
OnErr:
	if err != nil {
		option.OnError(err)
		return
	}
	// 获取登录 sign
	resp, err := httpClient.Do(newRequest("GET", "https://passport.baidu.com/v2/api/getqrcode?lp=pc").WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	body := readHtml(resp.Body)
	//log.Println(body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		goto OnErr
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取登录二维码错误, code:" + fmt.Sprint(data["errno"]))
	}
	sign := data["sign"].(string)
	img := "https://" + data["imgurl"].(string)
	link := fmt.Sprintf("https://wappass.baidu.com/wp/?qrlogin&t=%d"+
		"&error=0&sign=%s&cmd=login&lp=pc&tpl=netdisk&uaonly="+
		"&client_id=&adapter=3&client=&qrloginfrom=pc&wechat=0&traceid=", time.Now().Unix(), sign)
	if err != nil {
		goto OnErr
	}
	option.OnQrCode(img, link)
	//log.Println(link)
	err, uss := func() (err error, uss string) {
		for {
			var resp *http.Response
			resp, err = httpClient.Do(newRequest("GET",
				fmt.Sprintf("https://passport.baidu.com/channel/unicast?channel_id=%s"+
					"&tpl=netdisk&callback=&apiver=v3&tt=%d&_=%d", sign,
					time.Now().UnixNano()/int64(time.Millisecond), time.Now().UnixNano()/int64(time.Millisecond))).
				WithContext(option.Ctx))
			if err != nil {
				return
			}
			str := readHtml(resp.Body)
			str = strings.Trim(str, "() \n")
			log.Println(str)
			var ret tJson
			err = json.Unmarshal(tBin(str), &ret)
			if err != nil {
				return
			}
			if ret["errno"].(float64) == 0 {
				var c tJson
				err = json.Unmarshal(tBin(ret["channel_v"].(string)), &c)
				if err != nil {
					return
				}
				if c["status"].(float64) == 0 {
					uss = c["v"].(string)
					return
				}
			}
		}
	}()
	if err != nil {
		goto OnErr
	}
	link = fmt.Sprintf("https://passport.baidu.com/v3/login/main/qrbdusslogin?"+
		"v=%d&bduss=%s&u=&loginVersion=v4&qrcode=1&tpl=netdisk&apiver=v3"+
		"&tt=%d&traceid=&callback=bd__cbs__cupstt",
		time.Now().UnixNano()/int64(time.Millisecond), uss, time.Now().UnixNano()/int64(time.Millisecond))
	log.Println(link)
	resp, err = httpClient.Do(newRequest("GET", link))
	if err != nil {
		goto OnErr
	}
	body = readHtml(resp.Body)
	if strings.Contains(body, `"no": "0"`) {
		handleLoginSuccess()
		option.OnSuccess()
		return
	}
	err = errors.New(body)
	goto OnErr
}

func BaiduCookieLogin(cookieStr string) error {
	cookieJar, _ = cookiejar.New(nil)
	httpClient.Jar = cookieJar
	var cookies []*http.Cookie
	parts := strings.Split(strings.TrimSpace(cookieStr), ";")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
		if len(parts[i]) == 0 {
			continue
		}
		name, val := parts[i], ""
		if j := strings.Index(name, "="); j >= 0 {
			name, val = name[:j], name[j+1:]
		}
		cookies = append(cookies, &http.Cookie{Name: name, Value: val, Domain: ".baidu.com"})
	}

	u, _ := url.Parse("https://pan.baidu.com")
	cookieJar.SetCookies(u, cookies)
	log.Println(cookieJar.Cookies(u))
	return handleLoginSuccess()
}
