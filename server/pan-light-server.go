package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := iris.New()
	app.Get("/", func(ctx context.Context) {
		ctx.Write([]byte("Hello pan-light"))
	})
	app.Run(iris.Addr(":8081"))
}
