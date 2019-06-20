package pc_api

import (
	"github.com/kataras/iris"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/pc-api/middleware"
)

func Init(app *iris.Application) {
	app.Post("api/pc/login-token", artisan.ApiHandler(handleLoginToken))
	app.Post("api/pc/login", artisan.ApiHandler(handleLogin))

	pc := app.Party("api/pc/")
	pc.Use(middleware.PcJwtAuth)
	pc.Done()
	pc.Post("feedback", artisan.ApiHandler(handleFeedBack))
}
