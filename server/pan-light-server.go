package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/demo"
	"github.com/peterq/pan-light/server/pc-api"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := artisan.NewApp()
	//app.Get("/", func(ctx context.Context) {
	//	ctx.Write([]byte("Hello pan-light"))
	//})
	app.Get("/", func(ctx context.Context) {})

	cnf, ok := os.LookupEnv("pan_light_server_conf")
	if !ok {
		panic("the conf path must be specified")
	}
	configuration := iris.YAML(cnf)
	demo.Init(app.Party("/demo"), configuration.Other["demo"].(map[interface{}]interface{}))

	app.Use(artisan.ApiRecover)
	app.StaticWeb("/", "./static")
	pc_api.Init(app)
	app.Run(iris.Addr("127.0.0.1:8081"))
}
