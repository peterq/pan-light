package login

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type WxLoginOption struct {
	Ctx         context.Context
	OnError     func(err error)
	OnQrCodeUrl func(url string)
	OnScan      func()
	OnConfirm   func()
	OnSuccess   func()
}

func WxLogin(option *WxLoginOption) {
	var err error
OnErr:
	if err != nil {
		option.OnError(err)
		return
	}
	// 得到微信页面链接
	resp, err := httpClient.Do(newRequest("GET", "gotoWx").WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	link := getWxPageLink(readHtml(resp.Body))
	log.Println(link)
	// 解析回调地址
	u, err := url.Parse(link)
	if err != nil {
		goto OnErr
	}
	callbackLink := u.Query().Get("redirect_uri")
	log.Println(callbackLink)
	// 得到微信登录uid
	resp, err = httpClient.Do(newRequest("GET", link).WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	uid := getWxQrCodeUid(readHtml(resp.Body))
	log.Println(uid)
	option.OnQrCodeUrl("https://open.weixin.qq.com/connect/qrcode/" + uid)
	// 轮训扫描状态得到code
	wxCode, err := func() (code string, err error) {
		last := 408
		for {
			var resp *http.Response
			resp, err = httpClient.Do(newRequest("GET",
				fmt.Sprintf("https://long.open.weixin.qq.com/connect"+
					"/l/qrconnect?uuid=%s&_=%d&last=%d",
					uid, time.Now().UnixNano()/int64(time.Millisecond), last)).WithContext(option.Ctx))
			if err != nil {
				return
			}
			str := readHtml(resp.Body)
			log.Println(str)
			if strings.Index(str, "wx_errcode=408;") > 0 {
				last = 408
			} else if strings.Index(str, "wx_errcode=402;") > 0 {
				err = errors.New("超时, 请重试")
			} else if strings.Index(str, "wx_errcode=404;") > 0 {
				option.OnScan()
				last = 404
			} else if strings.Index(str, "wx_errcode=405;") > 0 {
				option.OnConfirm()
				reg := regexp.MustCompile(`window\.wx_code='(.*)'`)
				find := reg.FindStringSubmatch(str)
				code = find[1]
				return
			} else {
				err = errors.New("未知错误")
				break
			}
		}
		return
	}()
	if err != nil {
		goto OnErr
	}
	// 回调百度
	resp, err = httpClient.Do(newRequest("GET", callbackLink+"&code="+wxCode).WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	body := readHtml(resp.Body)
	reg := regexp.MustCompile(`next_url: (".+"),`)
	find := reg.FindStringSubmatch(body)
	if len(find) != 2 {
		err = errors.New("未知错误")
		log.Println(body)
		goto OnErr
	}
	link = ""
	json.Unmarshal([]byte(find[1]), &link)
	log.Println(link)
	if strings.Contains(link, "/account/bind") {
		err = errors.New("该微信号未绑定百度账号, 请绑定后再试")
		log.Println(cookieJar)
		goto OnErr
	} else if strings.Contains(link, "https://pan.baidu.com/disk/home") {
		handleLoginSuccess()
		option.OnSuccess()
	} else {
		err = errors.New("未知错误: " + link)
		goto OnErr
	}
}

var wxPageLinkReg = regexp.MustCompile(`"(https://open.weixin.qq.com/connect/qrconnect.*)"`)

func getWxPageLink(html string) string {
	str := wxPageLinkReg.FindStringSubmatch(html)
	return str[1]
}

var wxQrCOdeReg = regexp.MustCompile(`"/connect/qrcode/(.*)"`)

func getWxQrCodeUid(html string) string {
	str := wxQrCOdeReg.FindStringSubmatch(html)
	return str[1]
}
