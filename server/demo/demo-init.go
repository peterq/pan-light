package demo

import (
	"github.com/kataras/iris"
	"github.com/peterq/pan-light/server/demo/rpc"
)

func Init(router iris.Party, conf map[interface{}]interface{}) {
	// 静态页
	router.StaticWeb("/", "../static/noVNC")
	hosts := map[string]string{}
	for _, host := range conf["hosts"].([]interface{}) {
		hosts[host.(map[interface{}]interface{})["name"].(string)] = host.(map[interface{}]interface{})["password"].(string)
	}
	rpc.Init(router, hosts)
}
