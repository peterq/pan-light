package pc_api

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/pc-api/middleware"
	"time"
)

func Init(app *iris.Application) {
	app.Post("api/pc/login-token", artisan.ApiHandler(handleLoginToken))
	app.Post("api/pc/login", artisan.ApiHandler(handleLogin))

	pc := app.Party("api/pc/")
	pc.Use(middleware.PcJwtAuth)
	pc.Use(artisan.Throttle(artisan.ThrottleOption{ // 防止攻击
		Duration: time.Second * 5,
		Number:   20,
		GetKey: func(ctx context.Context) string {
			return "pc.api." + middleware.ContextLoginInfo(ctx).Uk()
		},
	}))
	pcAuthRoutes(pc)
}

// 需要登录的api
func pcAuthRoutes(r router.Party) {

	post := func(path string, handlers ...interface{}) {
		var hs []context.Handler
		for _, h := range handlers {
			if fn, ok := h.(func(ctx context.Context, param artisan.JsonMap) (result interface{}, err error)); ok {
				hs = append(hs, artisan.ApiHandler(fn))
			} else if o, ok := h.(artisan.ThrottleOption); ok {
				hs = append(hs, artisan.Throttle(o))
			} else {
				hs = append(hs, h.(context.Handler))
			}
		}
		r.Post(path, hs...)
	}

	post("feedback", artisan.ThrottleOption{
		Duration: time.Hour,
		Number:   5,
	}, handleFeedBack)

	post("refresh-token", handleRefreshToken)

	post("share", artisan.ThrottleOption{
		Duration: time.Hour,
		Number:   5,
	}, handleShareToSquare)

	post("share/list", artisan.ThrottleOption{
		Duration: time.Second * 5,
		Number:   5,
	}, handleShareList)

	post("share/hit", artisan.ThrottleOption{
		Duration: time.Second * 5,
		Number:   5,
	}, handleHitShare)

	post("link/md5", artisan.ThrottleOption{
		Duration: time.Second * 5,
		Number:   5,
	}, handleLinkMd5)
}
