package functions

import (
	"fmt"
	"github.com/peterq/pan-light/pc/pan-download"
	"os"
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
		useVip := false
		if u, ok := p["useVip"]; ok {
			useVip = u.(bool)
		}
		err := pan_download.Resume(p["downloadId"].(string), p["bin"].(string), useVip)
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
		if p["deleteFile"].(bool) {
			os.Remove(p["path"].(string))
		}
		return pan_download.Delete(p["downloadId"].(string))
	},
}

var downloadAsyncRoutes = map[string]asyncHandler{}
