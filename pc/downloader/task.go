package downloader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/peterq/pan-light/pc/downloader/internal"
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

type TaskState string

const (
	WaitStart   TaskState = "wait.start"
	WaitResume            = "wait.resume"
	COMPLETED             = "completed"
	STARTING              = "starting"
	DOWNLOADING           = "downloading"
	PAUSING               = "pausing"
	ERRORED               = "errored"
)

var noMoreSeg = errors.New("no more seg") // 所有seg分配完毕

type Task struct {
	id               TaskId
	fileId           string // 文件标识
	manager          *Manager
	linkResolver     LinkResolver
	requestDecorator func(*http.Request) *http.Request
	coroutineNumber  int
	segmentSize      int64
	savePath         string // 保存地址
	httpClient       *http.Client

	state         TaskState // 任务当前状态
	lastErr       error     // 保存上次错误
	link          string    // 链接地址
	finalLink     string    // redirect 之后的地址
	downloadCount int64
	speedCount    int64
	speed         int64

	fileLength        int64      // 文件总大小
	undistributed     []*segment // 尚未分配的片段
	wroteToDisk       []*segment // 文件内容写入磁盘的情况
	undistributedLock sync.Mutex
	wroteToDiskLock   sync.Mutex
	lastCaptureTime   time.Time // 上次快照时间

	workers               map[int]*worker // 工作协程map
	workersLock           sync.Mutex
	fileHandle            *os.File
	cancelSpeedCoroutine  context.CancelFunc
	speedCoroutineContext context.Context
}

func (task *Task) Id() TaskId {
	return task.id
}

func (task *Task) pause() error {
	if task.state != DOWNLOADING {
		return errors.New("当前状态不能暂停任务")
	}
	task.updateState(PAUSING)
	for _, w := range task.workers {
		w.cancel()
	}
	return nil
}

// 初始化生成下载状态
func (task *Task) init() (err error) {
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
	if !supportRange {
		return errors.New("该文件不支持并行下载")
	}
	task.fileHandle, err = os.OpenFile(task.savePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "打开本地文件错误")
	}
	task.undistributed = append(task.undistributed, &segment{
		start:  0,
		len:    task.fileLength,
		finish: 0,
	})
	task.downloadCount = 0
	for _, seg := range task.wroteToDisk {
		removeSegment(task.undistributed, seg)
		task.downloadCount += seg.len
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

// 获取当前速度
func (task *Task) getSpeed() int64 {
	return atomic.LoadInt64(&task.speed)
}

// 开始一个任务
func (task *Task) start() (err error) {
	if task.state != WaitStart {
		return errors.New("当前状态不能开始任务")
	}
	task.updateState(STARTING)
	err = task.init()
	if err != nil {
		task.lastErr = err
		task.updateState(ERRORED)
		return errors.Wrap(err, "任务初始化出错")
	}
	task.updateState(DOWNLOADING)
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
	return seg, nil
}

// 写入数据到磁盘
func (task *Task) writeToDisk(from int64, buffer *bytes.Buffer) (err error) {
	task.wroteToDiskLock.Lock()
	defer task.wroteToDiskLock.Unlock()
	_, err = task.fileHandle.Seek(from, io.SeekStart)
	if err != nil {
		return errors.Wrap(err, "文件seek错误")
	}
	l, err := buffer.WriteTo(task.fileHandle)
	//log.Println("写入片段", from, l)
	if err != nil {
		return errors.Wrap(err, "文件写入错误")
	}
	putBackSegment(task.wroteToDisk, &segment{
		start:  from,
		len:    l,
		finish: l,
	})
	task.capture()
	return
}

func (task *Task) capture() {
	if time.Now().Sub(task.lastCaptureTime) < time.Second {
		return
	}
	task.lastCaptureTime = time.Now()
	c := &internal.TaskCapture{
		Fid:       string(task.id),
		SavePath:  task.savePath,
		Completed: []*internal.FinishSeg{},
		Length:    task.fileLength,
	}
	for _, seg := range task.wroteToDisk {
		c.Completed = append(c.Completed, &internal.FinishSeg{
			Start: seg.start,
			Len:   seg.len,
		})
	}
	bin, err := proto.Marshal(c)
	if err != nil {
		log.Println("快照编码错误", err)
		return
	}
	task.notifyEvent("capture", string(bin))
}

// 下载出错, 放回片段到未下载
func (task *Task) downloadSegmentError(seg *segment) {
	task.undistributedLock.Lock()
	defer task.undistributedLock.Unlock()
	seg.finish = 0
	task.undistributed = putBackSegment(task.undistributed, seg)
	log.Println("下载片段错误", seg)
	logErr(fmt.Sprint(seg))
}

func logErr(strContent string) {
	fd, _ := os.OpenFile("seg.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fdTime := time.Now().Format("2006-01-02 15:04:05")
	fdContent := strings.Join([]string{"======", fdTime, "=====", strContent, "\n"}, "")
	buf := []byte(fdContent)
	fd.Write(buf)
	fd.Close()
}

// 下载成功, 放回片段到已下载
func (task *Task) downloadSegmentSuccess(seg *segment) {
	if seg.len == seg.finish {
		return
	}
	seg2 := &segment{
		start:  seg.start + seg.finish,
		len:    seg.len - seg.finish,
		finish: 0,
	}
	task.undistributedLock.Lock()
	defer task.undistributedLock.Unlock()
	task.undistributed = putBackSegment(task.undistributed, seg2)
}

// 当有工作线程退出时的回调
func (task *Task) onWorkerExit(w *worker) {
	task.workersLock.Lock()
	defer task.workersLock.Unlock()
	delete(task.workers, w.id)
	log.Println(fmt.Sprintf("task %s, worker %d exit", task.id, w.id))
	if len(task.workers) == 0 {
		go task.onAllWorkerExit()
	}
}

// 当所有线程退出时的回调
func (task *Task) onAllWorkerExit() {
	log.Println("所有worker结束")
	task.cancelSpeedCoroutine()
	task.fileHandle.Close()
	st := WaitStart
	if len(task.undistributed) == 0 {
		st = COMPLETED
	}
	task.updateState(st)
	log.Println(task.undistributed)
}

// 通知事件给外部
func (task *Task) notifyEvent(event string, data interface{}) {
	task.manager.eventNotify(&DownloadEvent{
		TaskId: task.id,
		Event:  event,
		Data:   data,
	})
}

// 更新任务状态
func (task *Task) updateState(state TaskState) {
	task.state = state
	data := map[string]interface{}{
		"state": state,
	}
	if state == ERRORED {
		data["error"] = task.lastErr.Error()
	}
	task.notifyEvent("task.state", data)
}

// 恢复任务
func (task *Task) resume(bin []byte) (err error) {
	if task.state != WaitResume {
		return errors.New("任务当前状态不能resume")
	}
	defer func() {
		if err != nil {
			task.lastErr = errors.Wrap(err, "任务恢复出错")
			task.updateState(ERRORED)
		}
	}()
	var data internal.TaskCapture
	err = proto.Unmarshal(bin, &data)
	if err != nil {
		return errors.Wrap(err, "无法decode数据")
	}
	task.fileId = data.Fid
	task.savePath = data.SavePath
	task.fileLength = data.Length
	for _, seg := range data.Completed {
		putBackSegment(task.wroteToDisk, &segment{
			start:  seg.Start,
			len:    seg.Len,
			finish: seg.Len,
		})
	}
	task.updateState(WaitStart)
	return nil
}

// 把一个段放回到一个slice中, 并进行必要的合并
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

// seg做减法, 只应用于恢复任务时
func removeSegment(queue []*segment, seg *segment) []*segment {
	head := seg.start
	tail := seg.start + seg.len
	// 头部衔接
	for idx := 0; idx < len(queue); idx++ {
		segInQueue := queue[idx]
		if segInQueue.start <= head && segInQueue.start+segInQueue.len >= tail {
			// 完全重合, 直接去掉
			if segInQueue.start == head && segInQueue.len == seg.len {
				return append(queue[:idx], queue[idx+1:]...)
			}
			// 头部重合, 留下后半段
			if segInQueue.start == head {
				segInQueue.start += seg.len
				segInQueue.len -= seg.len
				return queue
			}
			// 尾部重合, 留下头部
			if segInQueue.start+segInQueue.len == tail {
				segInQueue.len -= seg.len
				return queue
			}
			// 包含其中, 拆分
			seg2 := &segment{
				start: tail + 1,
				len:   segInQueue.start + segInQueue.len - tail,
			}
			segInQueue.len = seg.start - segInQueue.start
			rear := append([]*segment{}, queue[idx:]...)
			return append(append(queue[:idx], seg2), rear...)
		}
	}
	log.Print("去除错误, 未找到包含的段")
	return queue
}
