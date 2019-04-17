package instance

import (
	"context"
	"github.com/peterq/pan-light/demo/realtime"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const dockerImage = "pan-light-slave"

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
	pid         int // 进程号
	viewPwd     string
	operatePwd  string

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
	// 删除已有容器
	defer exec.Command("docker", "rm", "-v", "-f", h.SlaveName).Run()
	exec.Command("docker", "rm", "-v", "-f", h.SlaveName).Run()
	// 启动docker
	dockerP := exec.Command("docker", "run",
		"-m", "200m", "--memory-swap", "400m", // 200m 内存
		"--cpu-period=100000", "--cpu-quota=20000", // 20% cpu
		"-e", "vnc_operate_pwd="+h.operatePwd, "-e", "vnc_view_pwd="+h.viewPwd, // vnc 密码
		"--name='"+h.SlaveName+"'", // 容器名
		dockerImage)
	dockerP.Stdout = os.Stdout
	dockerP.Stderr = os.Stderr
	dockerP.Start()
	h.pid = dockerP.Process.Pid
	// 查询ip
	bin, err := exec.Command("sh", "-c", "docker inspect --format '{{ .NetworkSettings.IPAddress }}' "+h.SlaveName).Output()
	if err != nil {
		log.Println("获取ip错误", err)
		return
	}
	addr := strings.Trim(string(bin), "\r\n") + ":5091"
	// 配置地址
	func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		h.vncAddr = addr
		h.vncAddrCond.Broadcast() // 通知等待链接的代理
	}()
	dockerP.Wait()
}

func (h *Holder) VncProxy(rw io.ReadWriteCloser, proxyCb func(err error), viewOnly bool) {
	var addr string
	h.vncAddr = "127.0.0.1:5901"
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

	ctx, cancel := context.WithCancel(context.Background())

	var messageSize = 1024
	var ReadLoop = func(d io.Reader) {
		for {
			buffer := make([]byte, messageSize)
			n, err := d.Read(buffer)
			log.Println("r", buffer[:n])
			if err != nil {
				log.Println("Datachannel closed; Exit the ReadLoop:", err)
				cancel()
				return
			}
			con.Write(buffer[:n])
		}
	}

	var WriteLoop = func(d io.Writer) {
		for {
			buffer := make([]byte, messageSize)
			n, err := con.Read(buffer)
			log.Println("w", buffer[:n])
			if err != nil {
				log.Println("Datachannel closed; Exit the WriteLoop:", err)
				cancel()
				return
			}
			d.Write(buffer[:n])
		}
	}

	log.Println("proxy for rw", rw)
	go ReadLoop(rw)
	go WriteLoop(rw)
	<-ctx.Done()
	log.Println("proxy gone rw", rw)
}

func (h *Holder) VncProxyForOperation(rw io.ReadWriteCloser, proxyCb func(err error), order int64, ticket string) {
	if order != h.order || ticket != h.ticket {
		proxyCb(errors.New("ticket 验证失败"))
		return
	}
	h.VncProxy(rw, proxyCb, false)
}
