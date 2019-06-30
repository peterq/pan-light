package pan_download

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/downloader"
	"github.com/peterq/pan-light/pc/pan-api"
	"github.com/peterq/pan-light/pc/server-api"
	"github.com/peterq/pan-light/pc/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var manager *downloader.Manager

func init() {
	dep.OnInit(func() {
		parallel := 1024
		manager = &downloader.Manager{
			CoroutineNumber:       32,
			SegmentSize:           1024 * 1024 * 2,
			WroteToDiskBufferSize: 1024 * 512,
			LinkResolver:          LinkResolver,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					MaxIdleConns:    parallel,
					MaxConnsPerHost: parallel,
				},
			},
		}
		manager.Init()
		go func() {
			for evt := range manager.EventChan {
				go handleDownloadEvent(evt)
			}
		}()
		//go test()
	})
}

func Manager() *downloader.Manager {
	return manager
}

type linkTime struct {
	link string
	time time.Time
}

func (l *linkTime) expired() bool {
	return time.Now().Sub(l.time) > time.Hour
}

var linkCacheMap = map[string]linkTime{}

func LinkResolver(fileId string) (link string, err error) {
	if c, ok := linkCacheMap[fileId]; ok {
		if !c.expired() {
			return c.link, nil
		}
	}
	defer func() {
		if err == nil && link != "" {
			linkCacheMap[fileId] = linkTime{
				link: link,
				time: time.Now(),
			}
		}
	}()
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("链接解析严重错误: " + fmt.Sprint(e))
		}
	}()
	log.Println("resolve", fileId)
	args := strings.Split(fileId, ".")
	switch args[0] {
	case "vip":
		return vipLink(args[1])
	case "direct":
		return pan_api.LinkDirect(args[1])
	case "share":
		fileSize, _ := strconv.ParseInt(args[3], 10, 64)
		return VipLinkByMd5(args[1], args[2], fileSize)
	case "link":
		return decodeHyperLink(args[1])
	default:
		err = errors.New("unknown download method: " + args[0])
	}
	return
}
func decodeHyperLink(s string) (string, error) {
	bin, err := base64.StdEncoding.DecodeString(s)
	return string(bin), err
}

func encodeHyperLink(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func vipLink(fileId string) (link string, err error) {
	md5, sliceMd5, fileSize, err := RapidUploadMd5(fileId)
	if err != nil {
		err = errors.Wrap(err, "获取文件md5错误")
		return
	}
	return VipLinkByMd5(md5, sliceMd5, fileSize)
}

func handleDownloadEvent(event *downloader.DownloadEvent) {
	dep.NotifyQml("task.event", map[string]interface{}{
		"type":   event.Event,
		"taskId": event.TaskId,
		"data":   event.Data,
	})
}

func test() {
	//fileCompare()
	//return
	time.Sleep(3 * time.Second)
	id, err := DownloadFile("direct.730136432970379", "./yx.mp4")
	//id, err := DownloadFile("835313540804", "./project.mp4")
	log.Println(id, err)
}

func DownloadFile(fid, savePath string) (taskId downloader.TaskId, err error) {
	//savePath, err = filepath.Abs(savePath)
	if err != nil {
		return
	}
	taskId, err = manager.NewTask(fid, savePath, requestDecorator)
	return
}

func RapidUploadMd5(fid string) (md5, sliceMd5 string, fileSize int64, err error) {
	link, err := pan_api.LinkDirect(fid)
	if err != nil {
		err = errors.Wrap(err, "解析直链错误")
		return
	}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		errors.Wrap(err, "无法创建request")
		return
	}
	requestDecorator(req)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, 256*1024-1))
	resp, err := manager.HttpClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "访问直链错误")
		return
	}
	md5 = resp.Header.Get("Content-Md5")
	s := resp.Header.Get("Content-Range")
	s = strings.Trim(s, "]")
	s = strings.Split(s, "/")[1]
	fileSize, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "获取文件大小失败")
		return
	}
	defer resp.Body.Close()
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "获取前256k内容错误")
		return
	}
	if len(bin) != 256*1024 {
		err = errors.New("文件内容小于256k")
	}
	sliceMd5 = util.Md5bin(bin)
	return
}

func VipLinkByMd5(md5, sliceMd5 string, fileSize int64) (link string, err error) {
	result, err := server_api.Call("link-md5", map[string]interface{}{
		"md5":      md5,
		"sliceMd5": sliceMd5,
		"fileSize": fileSize,
	})
	if err != nil {
		err = errors.Wrap(err, "调用vip链接接口错误")
		return
	}
	link = result.(string)
	return
}

func requestDecorator(request *http.Request) *http.Request {
	request.Header.Set("User-Agent", pan_api.BaiduUA)
	return request
}

func Resume(id string, bin string, useVip bool) error {
	return manager.Resume(map[downloader.TaskId]string{
		downloader.TaskId(id): bin,
	}, requestDecorator)
}

func State(id string) interface{} {
	return manager.State(downloader.TaskId(id))
}

func Start(id string) error {
	return manager.StartTask(downloader.TaskId(id))
}

func Pause(id string) error {
	return manager.PauseTask(downloader.TaskId(id))
}

func Delete(id string) error {
	return manager.CancelTask(downloader.TaskId(id))
}

func Progress(id string) int64 {
	return manager.Progress(downloader.TaskId(id))
}

func fileCompare() {
	f1, err := os.OpenFile("/home/peterq/dev/projects/go/github.com/peterq/pan-light/pc/yx.mp4", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("err", err)
		return
	}

	f2, err := os.OpenFile("/home/peterq/dev/projects/go/github.com/peterq/pan-light/pc/yx.ok.mp4", os.O_RDONLY, 0644)
	if err != nil {
		log.Println("err", err)
		return
	}
	s1 := make([]byte, 512*1024)
	s2 := make([]byte, 512*1024)
	from := 0
	for {
		n1, err := f1.Read(s1)
		if err != nil {
			log.Println("err", err)
			return
		}

		n2, err := f2.Read(s2)
		if err != nil {
			log.Println("err", err)
			return
		}
		if n1 != n2 {
			log.Println("n1, n2", n1, n2)
			return
		}
		if !bytes.Equal(s1, s2) {
			log.Println(from)
			ioutil.WriteFile("cmp", []byte(fmt.Sprintf("%v\n%v", s1, s2)), os.ModePerm)
			//return
		}
		from += n1
	}
	log.Println("----------------------end")
}
