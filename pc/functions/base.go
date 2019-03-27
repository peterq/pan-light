package functions

import (
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/storage"
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
	// 存数据
	"storage.set": func(p map[string]interface{}) (result interface{}) {
		storage.UserStorageSet(p["k"].(string), p["v"].(string))
		return
	},
	// 取数据
	"storage.get": func(p map[string]interface{}) (result interface{}) {
		return storage.UserStorageGet(p["k"].(string))
	},
}

var baseAsyncRoutes = map[string]asyncHandler{}
