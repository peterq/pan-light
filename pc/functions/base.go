package functions

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/pan-download"
	"github.com/peterq/pan-light/pc/storage"
	"os"
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
		dep.DoClose()
		os.Exit(2)
		return
	},
	// config
	"config": func(p map[string]interface{}) (result interface{}) {
		maxParallelCorutineNumber := int(p["maxParallelCorutineNumber"].(float64))
		pan_download.Manager().CoroutineNumber = maxParallelCorutineNumber
		return true
	},
}

var baseAsyncRoutes = map[string]asyncHandler{}
