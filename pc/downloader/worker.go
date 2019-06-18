package downloader

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

// 实际下载协程
type worker struct {
	id     int
	task   *Task
	cancel func()
	ctx    context.Context
}

func (w *worker) work() {
	w.ctx, w.cancel = context.WithCancel(context.Background())
	errorNumber := 0
	maxErrorNumber := 2
	// 循环下载片段
WorkLoop:
	for {

		// 检查是否有连续错误
		if errorNumber > maxErrorNumber {
			log.Println("too many errors occurred")
			break
		}

		// 判断是否被取消
		select {
		case <-w.ctx.Done():
			break WorkLoop
		default:
		}

		// 获取新的下载片段
		seg, err := w.task.distributeSegment()
		// 没有新的下载片段, 退出下载循环
		if err == noMoreSeg {
			break
		}
		// 错误计数
		if err != nil {
			errorNumber++
			continue
		}
		// 下载该片段
		err = w.downloadSeg(seg)
		// 下载出错, 放回队列
		if err != nil {
			log.Println(err)
			w.task.downloadSegmentError(seg)
			errorNumber++
			continue
		}
		// 本次下载没出错, 错误计数置零
		errorNumber = 0
		w.task.downloadSegmentSuccess(seg)
	}
}

func (w *worker) downloadSeg(seg *segment) (err error) {
	// 构造请求
	req, _ := http.NewRequest("GET", w.task.finalLink, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", seg.start, seg.start+seg.len-1))
	// 用户回调, 可以修改请求
	req = w.task.requestDecorator(req)
	// 发送请求
	var resp *http.Response
	resp, err = w.task.httpClient.Do(req)
	if err != nil {
		return
	}
	// 读取数据流
	buf := w.task.manager.getBuffer()
	defer w.task.manager.releaseBuffer(buf)
	reader := resp.Body
	defer reader.Close()

	s := make([]byte, 1024)
	buffLeft := 0
	// 循环读取流
ReadStream:
	for {
		bin := s[:]
		buffLeft = buf.Cap() - buf.Len()
		if buffLeft < len(s) {
			bin = s[:buffLeft]
		}
		var l int
		// 读取流,with context
		select {
		case <-func() chan bool {
			ch := make(chan bool)
			go func() {
				l, err = reader.Read(bin)
				close(ch)
			}()
			return ch
		}():
		case <-w.ctx.Done():
			bufLen := int64(buf.Len())
			if bufLen > 0 {
				err = w.task.writeToDisk(seg.start+seg.finish, buf)
				if err == nil {
					seg.finish += bufLen
				}
			} else {
				err = errors.New("canceled")
			}
			break ReadStream
		}
		if l > 0 { // 有数据, 写入缓存
			//func() { // 投毒检测
			//	s := make([]byte, 1024)
			//	if bytes.Equal(bin[:l], s[:l]) {
			//		log.Println("检测到投毒", seg)
			//		log.Println(resp.StatusCode, resp.Header)
			//	}
			//}()
			buf.Write(bin[:l])
			w.task.addDownloadCount(int64(l))
			if buf.Len() == buf.Cap() || err == io.EOF { // 缓存满了, 或者流尾, 写入磁盘
				bufLen := int64(buf.Len())
				writeErr := w.task.writeToDisk(seg.start+seg.finish, buf)
				if writeErr != nil {
					err = writeErr
					break
				}
				buf.Reset()          // 重置缓冲区
				seg.finish += bufLen // 片段写入磁盘偏移量
			}
		}

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			break
		}
	}
	return
}
