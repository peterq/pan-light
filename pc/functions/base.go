package functions

import "github.com/peterq/pan-light/pc/dep"

func init() {
	syncMap(baseSyncRoutes)
	asyncMap(baseAsyncRoutes)
}

var baseSyncRoutes = map[string]syncHandler{
	// 获取当前时间
	"env.internal_server_url": func(p map[string]interface{}) interface{} {
		return dep.Env.InternalServerUrl
	},
}

var baseAsyncRoutes = map[string]asyncHandler{}
