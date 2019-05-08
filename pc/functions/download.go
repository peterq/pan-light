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
		err := pan_download.Resume(p["downloadId"].(string), p["bin"].(string))
		if err != nil {
			return err
		}
		return true
	},
	"download.state": func(p map[string]interface{}) interface{} {
		return pan_download.State(p["downloadId"].(string))
	},
	"download.start": func(p map[string]interface{}) interface{} {
		return pan_download.Start(p["downloadId"].(string))
	},
	"download.pause": func(p map[string]interface{}) interface{} {
		return pan_download.Pause(p["downloadId"].(string))
	},
	"download.delete": func(p map[string]interface{}) interface{} {
		return pan_download.Delete(p["downloadId"].(string))
	},
}

var downloadAsyncRoutes = map[string]asyncHandler{}
