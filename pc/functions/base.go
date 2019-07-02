package functions

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/pan-api"
	"github.com/peterq/pan-light/pc/pan-download"
	"github.com/peterq/pan-light/pc/server-api"
	"github.com/peterq/pan-light/pc/storage"
	"log"
)

func init() {
	syncMap(baseSyncRoutes)
	asyncMap(baseAsyncRoutes)
}

var baseSyncRoutes = map[string]syncHandler{
	// 获取当前时间
	"env.internal_server_url": func(p map[string]interface{}) interface{} {
		return dep.Env.InternalServerUrl
	},
	// 版本
	"env.version": func(p map[string]interface{}) interface{} {
		return dep.Env.VersionString
	},
	// 存数据
	"storage.set": func(p map[string]interface{}) (result interface{}) {
		storage.UserStorageSet(p["k"].(string), p["v"].(string))
		return
	},
	// 取数据
	"storage.get": func(p map[string]interface{}) (result interface{}) {
		return storage.UserStorageGet(p["k"].(string))
	},
	// 重启
	"reboot": func(p map[string]interface{}) (result interface{}) {
		dep.Reboot()
		return
	},
	// 退出
	"exit": func(p map[string]interface{}) (result interface{}) {
		dep.DoClose()
		//os.Exit(0)
		return
	},
	// config
	"config": func(p map[string]interface{}) (result interface{}) {
		maxParallelCorutineNumber := int(p["maxParallelCorutineNumber"].(float64))
		pan_download.Manager().CoroutineNumber = maxParallelCorutineNumber
		return true
	},
	// 退出登录
	"logout": func(p map[string]interface{}) (result interface{}) {
		if p["remove"].(bool) {
			storage.UserState.Logout = true
		}
		storage.Global.CurrentUser = "default"
		dep.Reboot()
		return
	},
	// 账号列表
	"account.list": func(p map[string]interface{}) (result interface{}) {
		var accounts []string
		for key, v := range storage.Global.UserStateMap {
			if v.Logout || key == "default" || key == storage.UserState.Username {
				continue
			}
			accounts = append(accounts, key)
		}
		return accounts
	},
	// 切换账号
	"account.change": func(p map[string]interface{}) (result interface{}) {
		storage.Global.CurrentUser = p["username"].(string)
		dep.Reboot()
		return
	},
}

var baseAsyncRoutes = map[string]asyncHandler{
	"api.login": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		data, err := server_api.Call("login-token", gson{
			"uk": storage.UserState.Uk,
		})
		if err != nil {
			reject(err)
			return
		}
		token := data.(gson)["token"]
		filename := data.(gson)["filename"].(string)
		if err != nil {
			reject(err)
			return
		}
		fid, serverPath, err := pan_api.UploadText(token.(string), "auth."+filename)
		if err != nil {
			reject(err)
			return
		}
		defer pan_api.DeleteFile(serverPath)
		link, secret, err := pan_api.ShareFile(fid, "")
		log.Println(link, secret, err)
		if err != nil {
			reject(err)
			return
		}
		jwt, err := server_api.Call("login", gson{
			"link":   link,
			"secret": secret,
			"token":  token,
		})
		if err != nil {
			reject(err)
			return
		}
		storage.UserState.Token = jwt.(string)
		log.Println(jwt)
		resolve(jwt)
	},

	"api.call": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		data, err := server_api.Call(p["name"].(string), p["param"].(gson))
		if err != nil {
			if data == nil {
				data = gson{
					"success": false,
					"message": err.Error(),
					"code":    -1,
				}
			}
			reject(data)
			return
		}
		resolve(data)
	},
}
