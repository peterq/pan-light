package functions

import (
	"github.com/peterq/pan-light/pc/pan-download"
)

func init() {
	syncMap(downloadSyncRoutes)
	asyncMap(downloadAsyncRoutes)
}

var downloadSyncRoutes = map[string]syncHandler{
	"download.new": func(p map[string]interface{}) interface{} {
		taskId, err := pan_download.DownloadFile(p["fid"].(string), p["savePath"].(string))
		if err != nil {
			return err
		}
		return taskId
	},
}

var downloadAsyncRoutes = map[string]asyncHandler{}
