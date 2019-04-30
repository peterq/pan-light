package pan_download

import (
	"bytes"
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/downloader"
	"github.com/peterq/pan-light/pc/pan-api"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var manager *downloader.Manager

func init() {
	dep.OnInit(func() {
		parallel := 128
		manager = &downloader.Manager{
			CoroutineNumber:       parallel,
			SegmentSize:           1024 * 1024 * 2,
			WroteToDiskBufferSize: 1024 * 512,
			LinkResolver:          pan_api.Link,
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
		go test()
	})
}

func handleDownloadEvent(event *downloader.DownloadEvent) {
	if event.Event == "task.speed" {
		speed := float64(event.Data.(int64))
		log.Println(speed / 1024 / 1024)
	} else {
		//log.Println(event)
	}
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
	id, err := DownloadFile("730136432970379", "./yx.mp4")
	//id, err := DownloadFile("835313540804", "./project.mp4")
	log.Println(id, err)
}

func DownloadFile(fid, savePath string) (taskId downloader.TaskId, err error) {
	savePath, err = filepath.Abs(savePath)
	if err != nil {
		return
	}
	taskId, err = manager.NewTask(fid, savePath, requestDecorator)
	return
}

func requestDecorator(request *http.Request) *http.Request {
	request.Header.Set("User-Agent", pan_api.BaiduUA)
	return request
}

func Resume(id string, bin []byte) error {
	return manager.Resume(map[downloader.TaskId][]byte{
		downloader.TaskId(id): bin,
	}, requestDecorator)
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
