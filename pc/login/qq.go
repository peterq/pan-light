package login

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type QQLoginOption struct {
	Ctx         context.Context
	OnError     func(err error)
	OnQrCodeUrl func(url string)
	OnScan      func()
	OnConfirm   func()
	OnSuccess   func()
}

func QQLogin(option *QQLoginOption) {
	var err error
OnErr:
	if err != nil {
		option.OnError(err)
		return
	}
	// 得到qq页面链接
	resp, err := httpClient.Do(newRequest("GET", "gotoQQ").WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	link := getQQPageLink(readHtml(resp.Body))
	log.Println(link)
	// 解析回调地址
	u, err := url.Parse(link)
	if err != nil {
		goto OnErr
	}
	callbackLink := u.Query().Get("redirect_uri")
	state := u.Query().Get("state")
	log.Println(callbackLink)
	// 访问qq iframe 获取cookie
	_, err = httpClient.Do(
		newRequest("GET", "https://xui.ptlogin2.qq.com/cgi-bin/xlogin?appid=716027609&daid=383&style=33&login_text=%E6%8E%88%E6%9D%83%E5%B9%B6%E7%99%BB%E5%BD%95&hide_title_bar=1&hide_border=1&target=self&s_url=https%3A%2F%2Fgraph.qq.com%2Foauth2.0%2Flogin_jump&pt_3rd_aid=100312028&pt_feedback_link=http%3A%2F%2Fsupport.qq.com%2Fwrite.shtml%3Ffid%3D780%26SSTAG%3Dwww.baidu.com.appid100312028").
			WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	// 获取二维码(qrsig cookie)
	resp, err = httpClient.Do(newRequest("GET", "https://ssl.ptlogin2.qq.com/ptqrshow?appid=716027609&e=2&l=M&s=3&d=72&v=4&t=0.648199619005845&daid=383&pt_3rd_aid=100312028").
		WithContext(option.Ctx))
	if err != nil {
		goto OnErr
	}
	// base64编码图片
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto OnErr
	}
	str := base64.StdEncoding.EncodeToString(bin)
	option.OnQrCodeUrl("data:image/png;base64," + str)
	// 取cookie
	login_sig := ""
	qr_sig := ""
	for _, c := range httpClient.Jar.Cookies(newRequest("GET", "https://ssl.ptlogin2.qq.com/ptqrlogin").URL) {
		if c.Name == "pt_login_sig" {
			login_sig = c.Value
		} else if c.Name == "qrsig" {
			qr_sig = c.Value
		}
	}

	if login_sig == "" || qr_sig == "" {
		err = errors.New("轮询初始化失败")
		goto OnErr
	}
	// 轮询状态
	for {
		link = fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrlogin?u1=%s"+
			"&ptqrtoken=%d&ptredirect=0&h=1&t=1&g=1&from_ui=1"+
			"&ptlang=2052&action=0-0-%d"+
			"&js_ver=10291&js_type=1&login_sig=%s&pt_uistyle=40&aid=716027609&daid=383&pt_3rd_aid=100312028&",
			"https%3A%2F%2Fgraph.qq.com%2Foauth2.0%2Flogin_jump", hash33(qr_sig),
			time.Now().UnixNano()/int64(time.Millisecond), login_sig)
		resp, err = httpClient.Do(newRequest("GET", link))
		if err != nil {
			goto OnErr
		}
		body := readHtml(resp.Body)
		log.Println(body)

		if strings.Contains(body, "二维码未失效") {

		} else if strings.Contains(body, "二维码认证中") {
			option.OnScan()
		} else if strings.Contains(body, "登录成功") {
			err = handleQQLoginSuccess(body, callbackLink, state)
			if err != nil {
				goto OnErr
			}
			handleLoginSuccess()
			option.OnSuccess()
			return
		} else if strings.Contains(body, "二维码已失效") {
			err = errors.New("二维码已失效, 请重试")
			goto OnErr
		} else {
			err = errors.New("登录状态异常: " + body)
			goto OnErr
		}

		select {
		case <-time.After(3 * time.Second):
		case <-option.Ctx.Done():
			err = errors.New("cancel")
			goto OnErr
		}
	}
}

//function(t){for (var e = 0, i = 0, n = t.length; i<n; ++i)e+=(e<<5)+t.charCodeAt(i); return 2147483647&e}
func hash33(t string) int {
	e := 0
	for i, n := 0, len(t); i < n; i++ {
		e += (e << 5) + int(t[i])
	}
	return 2147483647 & e
}

var qqPageLinkReg = regexp.MustCompile(`"(https://graph.qq.com/oauth2.0/authorize.*)"`)

func getQQPageLink(html string) string {
	str := qqPageLinkReg.FindStringSubmatch(html)
	return str[1]
}

func handleQQLoginSuccess(js, callbackLink, state string) (err error) {
	reg := regexp.MustCompile(`'(https://ssl.ptlogin2.graph.qq.com/check_sig.*?)'`)
	link := reg.FindStringSubmatch(js)[1]
	resp, err := httpClient.Do(newRequest("GET", link))
	if !strings.Contains(readHtml(resp.Body), "qclogin_success") {
		err = errors.New("check sig 失败")
	}
	// post 到auth
	req := newRequest("POST", "https://graph.qq.com/oauth2.0/authorize")
	p_skey := ""
	for _, c := range httpClient.Jar.Cookies(req.URL) {
		if c.Name == "p_skey" {
			p_skey = c.Value
			break
		}
	}
	f := url.Values{}
	params := map[string]interface{}{
		"response_type": "code",
		"client_id":     "100312028",
		"redirect_uri":  callbackLink,
		"scope":         "get_user_info,add_share,get_other_info,get_fanslist,get_idollist,add_idol,get_simple_userinfo",
		"state":         state,
		"switch":        "",
		"from_ptlogin":  "1",
		"src":           "1",
		"update_auth":   "1",
		"openapi":       "80901010",
		"g_tk":          getToken(p_skey),
		"auth_time":     time.Now().UnixNano() / int64(time.Millisecond),
		"ui":            guid(),
	}
	for k, v := range params {
		f.Set(k, fmt.Sprint(v))
	}
	req, err = http.NewRequest("POST", "https://graph.qq.com/oauth2.0/authorize", strings.NewReader(f.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}
	reg = regexp.MustCompile(`next_url: (".+"),`)
	find := reg.FindStringSubmatch(readHtml(resp.Body))
	if len(find) != 2 {
		err = errors.New("授权失败")
		return
	}
	link = ""
	json.Unmarshal([]byte(find[1]), &link)
	log.Println(link)
	if strings.Contains(link, "/account/bind") {
		err = errors.New("该qq号未绑定百度账号, 请绑定后再试")
	} else if strings.Contains(link, "https://pan.baidu.com/disk/home") {
		// 成功
	} else {
		err = errors.New("未知错误: " + link)
	}
	return
}

func getToken(ps_key string) int {
	hash := 5381
	for i, l := 0, len(ps_key); i < l; i++ {
		hash += (hash << 5) + int(ps_key[i])
	}
	return hash & 0x7fffffff
}

func guid() string {
	tpl := "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"
	reg := regexp.MustCompile(`[xy]`)
	return reg.ReplaceAllStringFunc(tpl, func(c string) string {
		r := rand.Intn(16)
		v := 0
		if c == "x" {
			v = r
		} else {
			v = r&0x3 | 0x8
		}
		return hex.EncodeToString([]byte{byte(v)})
	})
}
