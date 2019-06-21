package functions

import (
	"context"
	"github.com/peterq/pan-light/pc/login"
	"log"
)

func init() {
	syncMap(loginSyncRoutes)
	asyncMap(loginAsyncRoutes)
}

var loginSyncRoutes = map[string]syncHandler{}

var loginAsyncRoutes = map[string]asyncHandler{
	// 微信扫码登录
	"login.wx": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		ctx, cancel := context.WithCancel(UiContext)
		go func() {
			for msg := range qmlMsg {
				if str, ok := msg.(string); ok {
					if str == "cancel" {
						cancel()
					}
				}
			}
		}()
		option := &login.WxLoginOption{
			Ctx: ctx,
			OnError: func(err error) {
				log.Println(err)
				reject(err.Error())
			},
			OnQrCodeUrl: func(url string) {
				progress(gson{
					"type": "qrCode",
					"url":  url,
				})
			},
			OnScan: func() {
				progress(gson{
					"type": "scan.ok",
				})
			},
			OnConfirm: func() {
				progress(gson{
					"type": "confirm",
				})
			},
			OnSuccess: func() {
				resolve("ok")
			},
		}
		login.WxLogin(option)
	},

	// qq扫码登录
	"login.qq": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		ctx, cancel := context.WithCancel(UiContext)
		go func() {
			for msg := range qmlMsg {
				if str, ok := msg.(string); ok {
					if str == "cancel" {
						cancel()
					}
				}
			}
		}()
		option := &login.QQLoginOption{
			Ctx: ctx,
			OnError: func(err error) {
				log.Println(err)
				reject(err.Error())
			},
			OnQrCodeUrl: func(url string) {
				progress(gson{
					"type": "qrCode",
					"url":  url,
				})
			},
			OnScan: func() {
				progress(gson{
					"type": "scan.ok",
				})
			},
			OnConfirm: func() {
				progress(gson{
					"type": "confirm",
				})
			},
			OnSuccess: func() {
				resolve("ok")
			},
		}
		login.QQLogin(option)
	},

	// 百度扫码登录
	"login.baidu": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		ctx, cancel := context.WithCancel(UiContext)
		go func() {
			for msg := range qmlMsg {
				if str, ok := msg.(string); ok {
					if str == "cancel" {
						cancel()
					}
				}
			}
		}()
		option := &login.BaiduLoginOption{
			Ctx: ctx,
			OnError: func(err error) {
				log.Println(err)
				reject(err.Error())
			},
			OnQrCode: func(img string, pageUrl string) {
				progress(gson{
					"type":    "qrCode",
					"img":     img,
					"pageUrl": pageUrl,
				})
			},
			OnScan: func() {
				progress(gson{
					"type": "scan.ok",
				})
			},
			OnConfirm: func() {
				progress(gson{
					"type": "confirm",
				})
			},
			OnSuccess: func() {
				resolve("ok")
			},
		}
		login.BaiduLogin(option)
	},
	"login.cookie": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		cookie := p["cookie"].(string)
		err := login.BaiduCookieLogin(cookie)
		if err != nil {
			reject(err)
		} else {
			resolve("ok")
		}
	},
}
