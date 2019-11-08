package demo

import (
	"github.com/kataras/iris/v12"
	"github.com/peterq/pan-light/server/demo/rpc"
)

func Init(router iris.Party, conf map[interface{}]interface{}) {
	// 静态页
	router.Get("/", func(context iris.Context) {
		context.Redirect("/demo/init.html")
	})
	hosts := map[string]string{}
	for _, host := range conf["hosts"].([]interface{}) {
		hosts[host.(map[interface{}]interface{})["name"].(string)] = host.(map[interface{}]interface{})["password"].(string)
	}
	rpc.Init(router, hosts)
}
