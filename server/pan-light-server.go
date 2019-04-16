package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/peterq/pan-light/server/demo"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := iris.New()
	app.Get("/", func(ctx context.Context) {})

	cnf, ok := os.LookupEnv("pan_light_server_conf")
	if !ok {
		panic("the conf path must be specified")
	}
	configuration := iris.YAML(cnf)
	demo.Init(app.Party("/demo"), configuration.Other["demo"].(map[interface{}]interface{}))

	app.Run(iris.Addr(":8081"))
}
