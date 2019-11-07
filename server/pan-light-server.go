package main

import (
	"github.com/kataras/iris/v12"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/demo"
	"github.com/peterq/pan-light/server/pc-api"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := artisan.NewApp()

	cnf, ok := os.LookupEnv("pan_light_server_conf")
	if !ok {
		panic("the conf path must be specified")
	}
	configuration := iris.YAML(cnf)
	demo.Init(app.Party("/demo"), configuration.Other["demo"].(map[interface{}]interface{}))

	app.Use(artisan.ApiRecover)
	app.HandleDir("/", "./static")
	pc_api.Init(app)
	app.Run(iris.Addr(":8081"), iris.WithConfiguration(iris.Configuration{
		DisablePathCorrectionRedirection: true,
		DisablePathCorrection:            true,
	}))
}
