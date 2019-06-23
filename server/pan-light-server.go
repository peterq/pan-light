package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/pc-api"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := artisan.NewApp()
	app.Get("/", func(ctx context.Context) {
		ctx.Write([]byte("Hello pan-light"))
	})
	app.Use(artisan.ApiRecover)
	pc_api.Init(app)
	app.Run(iris.Addr("127.0.0.1:8081"))
}
