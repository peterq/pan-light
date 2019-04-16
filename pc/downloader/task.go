package downloader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type TaskState int

const (
	WAITE_START TaskState = iota
	STARTING
	DOWNLOADING
	PAUSEING
	PAUSED
	ERRORED
)

var noMoreSeg = errors.New("no more seg") // 所有seg分配完毕

type Task struct {
	Id               TaskId
	fileId           string // 文件标识
	manager          *Manager
	linkResolver     LinkResolver
	requestDecorator func(*http.Request) *http.Request
	coroutineNumber  int
	segmentSize      int64
	savePath         string // 保存地址
	httpClient       *http.Client

	initialized   bool      // 是否初始化
	state         TaskState // 任务当前状态
	lastErr       error     // 保存上次错误
	link          string    // 链接地址
	finalLink     string    // redirect 之后的地址
	downloadCount int64
	speedCount    int64
	speed         int64

	fileLength        int64      // 文件总大小
	undistributed     []*segment // 尚未分配的片段
	distributed       []*segment // 已经分配的片段
	finished          []*segment // 已经完成的片段
	wroteToDisk       []*segment // 文件内容写入磁盘的情况
	distributeLock    sync.Mutex
	undistributedLock sync.Mutex
	distributedLock   sync.Mutex
	finishedLock      sync.Mutex
	wroteToDiskLock   sync.Mutex

	workers               map[int]*worker // 工作协程map
	workersLock           sync.Mutex
	fileHandle            *os.File
	fileLock              sync.Mutex
	cancelSpeedCoroutine  context.CancelFunc
	speedCoroutineContext context.Context
}

// 初始化生成下载状态
func (task *Task) init() (err error) {
	if task.initialized {
		return errors.New("重复初始化")
	}
	task.link, err = task.linkResolver(task.fileId)
	if err != nil {
		return errors.Wrap(err, "获取下载链接错误")
	}
	req, err := http.NewRequest("GET", task.link, nil)
	if err != nil {
		return errors.Wrap(err, "无法创建request")
	}
	req = task.requestDecorator(req)
	task.finalLink, err = redirectedLink(req)
	if err != nil {
		return errors.Wrap(err, "获取最终链接错误")
	}
	req.URL, _ = url.Parse(task.finalLink)
	var supportRange bool
	task.fileLength, _, supportRange, err = downloadFileInfo(req)
	if err != nil {
		return errors.Wrap(err, "获取文件信息错误")
	}
	task.undistributed = append(task.undistributed, &segment{
		start:  0,
		len:    task.fileLength,
		finish: 0,
		state:  segmentWait,
	})
	if !supportRange {
		return errors.New("该文件不支持并行下载")
	}
	task.fileHandle, err = os.OpenFile(task.savePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "打开本地文件错误")
	}
	task.workers = map[int]*worker{}
	return nil
}

func (task *Task) addDownloadCount(cnt int64) {
	atomic.AddInt64(&task.downloadCount, cnt)
	atomic.AddInt64(&task.speedCount, cnt)
}

func (task *Task) speedCalculateCoroutine() {
	t := time.Tick(time.Second)
Loop:
	for {
		select {
		case <-task.speedCoroutineContext.Done():
			atomic.SwapInt64(&task.speedCount, 0)
			break Loop
		case <-t:
			cnt := atomic.SwapInt64(&task.speedCount, 0)
			task.notifyEvent("task.speed", cnt)
			atomic.SwapInt64(&task.speed, cnt)
		}
	}
}

func (task *Task) getSpeed() int64 {
	return atomic.LoadInt64(&task.speed)
}

// 开始一个任务
func (task *Task) start() (err error) {
	if task.state != WAITE_START {
		return errors.New("当前状态不能开始任务")
	}
	err = task.init()
	if err != nil {
		task.state = ERRORED
		task.lastErr = err
		return errors.Wrap(err, "任务初始化出错")
	}
	task.speedCoroutineContext, task.cancelSpeedCoroutine = context.WithCancel(context.Background())
	go task.speedCalculateCoroutine()
	for i := 0; i < task.coroutineNumber; i++ {
		task.workers[i] = &worker{
			id:   i,
			task: task,
		}
		go func(w *worker) {
			w.work()
			task.onWorkerExit(w)
		}(task.workers[i])
	}
	return nil
}

// 分配一段下载任务
func (task *Task) distributeSegment() (seg *segment, err error) {
	task.undistributedLock.Lock()
	defer task.undistributedLock.Unlock()
	task.distributedLock.Lock()
	defer task.distributedLock.Unlock()
	segLen := len(task.undistributed)
	if segLen == 0 {
		return nil, noMoreSeg
	}
	seg = task.undistributed[segLen-1]
	// seg 过大, 拆分
	if seg.len > task.segmentSize*3/2 {
		seg1 := &segment{
			start: seg.start,
			len:   task.segmentSize,
		}
		seg2 := &segment{
			start: seg.start + seg1.len,
			len:   seg.len - seg1.len,
		}
		task.undistributed[segLen-1] = seg2
		seg = seg1
	} else {
		task.undistributed = task.undistributed[:segLen-1]
	}
	seg.state = segmentDownloading
	task.distributed = append(task.distributed, seg)
	return seg, nil
}

// 写入数据到磁盘
func (task *Task) writeToDisk(from int64, buffer *bytes.Buffer) (err error) {
	task.fileLock.Lock()
	defer task.fileLock.Unlock()
	_, err = task.fileHandle.Seek(from, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "文件seek错误")
	}
	_, err = buffer.WriteTo(task.fileHandle)
	//l, err := buffer.WriteTo(task.fileHandle)
	//log.Println("写入片段", from, l)
	if err != nil {
		return errors.Wrap(err, "文件写入错误")
	}
	return
}

// 下载出错, 放回片段到未下载
func (task *Task) downloadSegmentError(seg *segment) {
	task.undistributedLock.Lock()
	defer task.undistributedLock.Unlock()
	seg.state = segmentWait
	seg.finish = 0
	task.undistributed = putBackSegment(task.undistributed, seg)
	log.Println("下载片段错误", seg)
	logerr(fmt.Sprint(seg))
}

func logerr(str_content string) {
	fd, _ := os.OpenFile("seg.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{"======", fd_time, "=====", str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

// 下载成功, 放回片段到已下载
func (task *Task) downloadSegmentSuccess(seg *segment) {
	task.distributedLock.Lock()
	defer task.distributedLock.Unlock()
	task.finishedLock.Lock()
	defer task.finishedLock.Unlock()
	if seg.len == seg.finish {
		seg.state = segmentFinished
		putBackSegment(task.finished, seg)
		return
	}
	seg1 := &segment{
		start:  seg.start,
		len:    seg.finish,
		finish: seg.finish,
		state:  segmentFinished,
	}
	seg2 := &segment{
		start:  seg.start + seg.finish,
		len:    seg.len - seg.finish,
		finish: 0,
		state:  segmentWait,
	}
	task.finished = putBackSegment(task.finished, seg1)
	task.undistributed = putBackSegment(task.undistributed, seg2)
}

func (task *Task) onWorkerExit(w *worker) {
	task.workersLock.Lock()
	defer task.workersLock.Unlock()
	delete(task.workers, w.id)
	log.Println(fmt.Sprintf("task %d, worker %d exit", task.Id, w.id))
	if len(task.workers) == 0 {
		go task.onAllWorkerExit()
	}
}

func (task *Task) onAllWorkerExit() {
	log.Println("所有worker结束")
	task.cancelSpeedCoroutine()
	log.Println(task.undistributed)
	log.Println(task.distributed)
	log.Println(task.finished)
}

func (task *Task) notifyEvent(event string, data interface{}) {
	task.manager.eventNotify(&DownloadEvent{
		TaskId: task.Id,
		Event:  event,
		Data:   data,
	})
}

func putBackSegment(queue []*segment, seg *segment) []*segment {
	head := seg.start
	tail := seg.start + seg.len
	// 头部衔接
	for idx := 0; idx < len(queue); idx++ {
		segInQueue := queue[idx]
		if segInQueue.start+segInQueue.len+1 == head {
			if idx == len(queue)-1 {
				queue = queue[:idx]
			} else {
				queue = append(queue[:idx], queue[idx+1:]...)
			}
			segInQueue.len += seg.len
			seg = segInQueue
			break
		}
	}
	// 尾部衔接
	for idx := 0; idx < len(queue); idx++ {
		segInQueue := queue[idx]
		if segInQueue.start == tail+1 {
			if idx == len(queue)-1 {
				queue = queue[:idx]
			} else {
				queue = append(queue[:idx], queue[idx+1:]...)
			}
			seg.len += segInQueue.len
			break
		}
	}
	// 插入队列
	queue = append(queue, seg)
	return queue
}
