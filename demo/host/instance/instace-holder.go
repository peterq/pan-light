package instance

import (
	"context"
	"fmt"
	"github.com/peterq/pan-light/demo/realtime"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type gson = map[string]interface{}

type Holder struct {
	inited bool

	SlaveName    string
	HostName     string
	HostPassword string

	rt  *realtime.RealTime
	ctx context.Context

	checkUserChan chan struct{}

	vncAddr     string
	vncAddrLock sync.Mutex
	vncAddrCond *sync.Cond

	order     int64
	ticket    string
	sessionId int64
}

func (h *Holder) Init(rt *realtime.RealTime, ctx context.Context) {
	if h.inited {
		return
	}
	h.inited = true
	h.rt = rt
	h.ctx = ctx
	h.checkUserChan = make(chan struct{})
	h.vncAddrCond = sync.NewCond(&h.vncAddrLock)
	for {
		re, err := h.rt.Call("host.next.user", gson{
			"slave": h.SlaveName,
		})
		if err != nil {
			<-h.checkUserChan
			continue
		}
		result := re.(gson)
		h.order = int64(result["order"].(float64))
		h.ticket = result["ticket"].(string)
		h.sessionId = int64(result["sessionId"].(float64))
		h.startIns()
		h.order = -1
	}
}

func (h *Holder) CheckUser() {
	if h.checkUserChan != nil {
		select {
		case h.checkUserChan <- struct{}{}:
		default:
		}
	}
}

// 启动实例, 这个方法需要阻塞
func (h *Holder) startIns() {
	// 实例结束, 删除vnc地址
	defer func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		h.vncAddr = ""
	}()
	// 启动docker

	// 配置地址
	func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		h.vncAddr = "127.0.0.1:5901"
		h.vncAddrCond.Broadcast() // 通知等待链接的代理
	}()
	time.Sleep(10 * time.Minute)
}

func (h *Holder) VncProxy(rw io.ReadWriteCloser, proxyCb func(err error), viewOnly bool) {
	var addr string
	func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		for h.vncAddr == "" {
			h.vncAddrCond.Wait()
		}
		addr = h.vncAddr
	}()

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	con, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("connection failed: %v\n", err)
		proxyCb(errors.New("连接docker内部vnc出错"))
		return
	}
	defer con.Close()
	// 开始进行数据转发
	proxyCb(nil)

	var messageSize = 1024
	var ReadLoop = func(d io.Reader) {
		for {
			buffer := make([]byte, messageSize)
			n, err := d.Read(buffer)
			if err != nil {
				fmt.Println("Datachannel closed; Exit the readloop:", err)
				return
			}
			con.Write(buffer[:n])
		}
	}

	var WriteLoop = func(d io.Writer) {
		for {
			buffer := make([]byte, messageSize)
			n, err := con.Read(buffer)
			if err != nil {
				fmt.Println("Datachannel closed; Exit the readloop:", err)
				con.Close()
				return
			}
			d.Write(buffer[:n])
		}
	}

	go WriteLoop(rw)
	ReadLoop(rw)
}

func (h *Holder) VncProxyForOperation(rw io.ReadWriteCloser, proxyCb func(err error), order int64, ticket string) {
	if order != h.order || ticket != h.ticket {
		proxyCb(errors.New("ticket 验证失败"))
		return
	}
	h.VncProxy(rw, proxyCb, false)
}
