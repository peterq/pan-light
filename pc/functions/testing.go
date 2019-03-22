package functions

import (
	"io/ioutil"
	"net/http"
)

func init() {
	syncMap(testingSyncRoutes)
	asyncMap(testingAsyncRoutes)
}

var testingSyncRoutes = map[string]syncHandler{
	"add": func(p map[string]interface{}) interface{} {
		return p["a"].(float64) + p["b"].(float64)
	},
}

var testingAsyncRoutes = map[string]asyncHandler{
	"ip.info": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}),
		progress func(interface{}), qmlMsg chan interface{}) {
		r, e := http.Get("http://pv.sohu.com/cityjson?ie=utf-8")
		if e != nil {
			reject(e)
			return
		}
		bin, e := ioutil.ReadAll(r.Body)
		if e != nil {
			reject(e)
			return
		}
		resolve(string(bin))
	},
}
