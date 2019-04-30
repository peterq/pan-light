package functions

import (
	"fmt"
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
		return fmt.Sprint(taskId)
	},
	"download.resume": func(p map[string]interface{}) interface{} {
		err := pan_download.Resume(p["downloadId"].(string), []byte(p["bin"].(string)))
		if err != nil {
			return err
		}
		return true
	},
}

var downloadAsyncRoutes = map[string]asyncHandler{}
