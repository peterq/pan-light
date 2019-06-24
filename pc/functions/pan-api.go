package functions

import (
	"fmt"
	"github.com/peterq/pan-light/pc/pan-api"
	"github.com/peterq/pan-light/pc/pan-download"
)

func init() {
	syncMap(panApiSyncRoutes)
	asyncMap(panApiAsyncRoutes)
}

var panApiSyncRoutes = map[string]syncHandler{}

var panApiAsyncRoutes = map[string]asyncHandler{
	// pan api 初始化
	"pan.init": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		ctx, err := pan_api.GetSign()
		if err != nil {
			reject(err.Error())
		} else {
			resolve(ctx)
		}
	},

	"pan.ls": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		list, err := pan_api.ListDir(p["path"].(string))
		if err != nil {
			reject(err.Error())
		} else {
			resolve(list)
		}
	},

	"pan.link": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		link, err := pan_download.LinkResolver(p["fid"].(string))
		if err != nil {
			reject(err.Error())
		} else {
			resolve(link)
		}
	},
	"pan.usage": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		result, err := pan_api.Usage()
		if err != nil {
			reject(err.Error())
		} else {
			resolve(result)
		}
	},
	// 获取极速上传的md5
	"pan.rapid.md5": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		_, sliceMd5, _, err := pan_download.RapidUploadMd5(fmt.Sprint(int(p["fid"].(float64))))
		if err != nil {
			reject(err)
			return
		}
		resolve(sliceMd5)
	},
	// 通过md5转存文件
	"pan.save.md5": func(p map[string]interface{}, resolve func(interface{}), reject func(interface{}), progress func(interface{}), qmlMsg chan interface{}) {
		md5 := p["md5"].(string)
		sliceMd5 := p["sliceMd5"].(string)
		path := p["path"].(string)
		fileSize := int64(p["fileSize"].(float64))
		serverPath, _, _, err := pan_api.SaveFileByMd5(md5, sliceMd5, path, fileSize)
		if err != nil {
			reject(err)
			return
		}
		resolve(serverPath)
	},
}
