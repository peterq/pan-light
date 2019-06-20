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
			return "pc.api." + middleware.CotextLoginInfo(ctx).Uk()
		},
	}))
	pcAuthRoutes(pc)
}

// 需要登录的api
func pcAuthRoutes(r router.Party) {
	r.Post("feedback", artisan.Throttle(artisan.ThrottleOption{
		Duration: time.Hour,
		Number:   5,
	}), artisan.ApiHandler(handleFeedBack))
}
