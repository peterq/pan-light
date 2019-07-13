package instance

import (
	"context"
	"encoding/json"
	"github.com/peterq/pan-light/demo/realtime"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const dockerImage = "pan-light-slave"

type gson = map[string]interface{}

type notifier struct {
	cond *sync.Cond
	lock *sync.Mutex
}

func newNotifier() notifier {
	l := &sync.Mutex{}
	return notifier{
		cond: sync.NewCond(l),
		lock: l,
	}
}

func (n notifier) broadcast() {
	n.cond.Broadcast()
}

func (n notifier) wait() {
	n.lock.Lock()
	n.cond.Wait()
	n.lock.Unlock()
}

type Holder struct {
	inited bool

	SlaveName    string
	HostName     string
	HostPassword string
	WsAddr       string

	rt  *realtime.RealTime
	ctx context.Context

	checkUserChan   chan struct{} // 用于阻塞控制获取下一个体验用户
	slaveOkNotifier notifier

	vncAddr       string
	vncAddrLock   sync.Mutex
	vncAddrCond   *sync.Cond
	pid           int // 进程号
	viewPwd       string
	operatePwd    string
	containerName string

	order     int64
	ticket    string
	sessionId string
}

func (h *Holder) Order() int64 {
	return h.order
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
	h.containerName = h.SlaveName
	h.slaveOkNotifier = newNotifier()
	for {
		re, err := h.rt.Call("next.user", gson{
			"slave": h.SlaveName,
		})
		if err != nil {
			<-h.checkUserChan
			continue
		}
		result := re.(gson)
		h.order = int64(result["order"].(float64))
		h.ticket = result["ticket"].(string)
		h.sessionId = result["sessionId"].(string)
		h.viewPwd = "peter.q.is.so.cool"
		h.operatePwd = h.ticket
		h.startIns()
		h.rt.Emit("slave.exit", h.SlaveName)
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
	defer exec.Command("docker", "rm", "-v", "-f", h.containerName).Run()
	exec.Command("docker", "rm", "-v", "-f", h.containerName).Run()
	e, _ := filepath.Abs("./slave/ubuntu16.04/root.pan-light")
	log.Println(e)
	// 启动docker
	dockerP := exec.Command("docker", "run",
		"-m", "400m", "--memory-swap", "500m", // 400m 内存
		"--cpu-period=100000", "--cpu-quota=40000", // 40% cpu
		"-e", "vnc_operate_pwd="+h.operatePwd, "-e", "vnc_view_pwd="+h.viewPwd, // vnc 密码
		"-e", "host_name="+h.HostName, "-e", "host_password="+h.HostPassword, // host 密码
		"-e", "slave_name="+h.SlaveName, "-e", "ws_addr="+h.WsAddr, // ws 地址
		"-e", "demo_order="+strconv.FormatInt(h.order, 10), // demo order
		"-e", "demo_user="+h.sessionId, // 用户session
		//"-v"+e+":/root/pan-light",    // 开发时文件映射, 正式环境使用docker copy
		"--name="+h.containerName+"", // 容器名
		dockerImage)
	defer exec.Command("docker", "kill", h.containerName)
	dockerP.Stdout = os.Stdout
	dockerP.Stderr = os.Stderr
	dockerP.Start()
	defer dockerP.Process.Kill()
	h.pid = dockerP.Process.Pid
	// 查询ip
	go func() { // 最多等待35秒
		order := h.order
		<-time.After(35 * time.Second)
		if order != h.order {
			return
		}
		log.Println("docker run container failed")
		h.slaveOkNotifier.broadcast()
	}()
	h.slaveOkNotifier.wait()
	bin, err := exec.Command("docker", "inspect", "--format",
		"{{ .NetworkSettings.IPAddress }}", h.containerName).Output()
	if err != nil {
		log.Println("获取ip错误", err, bin)
		return
	}
	addr := strings.Trim(string(bin), "\r\n") + ":5901"
	log.Println("docker vnc addr", addr)
	// 配置地址
	func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		h.vncAddr = addr
		h.vncAddrCond.Broadcast() // 通知等待链接的代理
	}()
	// 超过6分强制结束进程
	exited := make(chan bool)
	defer close(exited)
	go func() {
		select {
		case <-time.After(6 * time.Minute):
			dockerP.Process.Kill()
		case <-exited:
		}
	}()
	dockerP.Wait()
}

func (h *Holder) VncProxy(rw io.ReadWriteCloser, proxyCb func(err error)) {
	defer rw.Close()
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

	ctx, cancel := context.WithCancel(context.Background())

	var messageSize = 1024
	var ReadLoop = func(d io.Reader) {
		for {
			buffer := make([]byte, messageSize)
			n, err := d.Read(buffer)
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

func (h *Holder) WsProxy(userConn *websocket.Conn) {
	defer userConn.Close()
	var addr string
	func() {
		h.vncAddrLock.Lock()
		defer h.vncAddrLock.Unlock()
		for h.vncAddr == "" {
			log.Println("ws proxy 等待docker ip")
			h.vncAddrCond.Wait()
		}
		addr = h.vncAddr
	}()

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", addr)
	vncConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("connection failed: %v\n", err)
		bin, _ := json.Marshal(gson{
			"proxyOk": false,
			"message": "连接docker内部vnc出错",
		})
		userConn.Write(bin)
		return
	}
	defer vncConn.Close()

	bin, _ := json.Marshal(gson{
		"proxyOk": true,
		"message": "连接成功",
	})
	userConn.Write(bin)

	ctx, cancel := context.WithCancel(context.Background())
	var ReadLoop = func() {
		for {
			var msg []byte
			err := websocket.Message.Receive(userConn, &msg)
			if err != nil {
				log.Println("ws agent read from user err:", err)
				cancel()
				return
			}
			vncConn.Write(msg)
		}
	}

	var WriteLoop = func() {
		for {
			msg := make([]byte, 1024)
			l, err := vncConn.Read(msg)
			if err != nil {
				log.Println("ws agent read from vnc err:", err)
				cancel()
				return
			}
			websocket.Message.Send(userConn, msg[:l])
		}
	}

	log.Println("ws proxy start", h.SlaveName)
	go ReadLoop()
	go WriteLoop()
	<-ctx.Done()
	log.Println("ws proxy gone", h.SlaveName)
}

// 处理 slave 发来的消息
func (h *Holder) HandleEvent(evt string, data interface{}) {
	if evt == "start.ok" {
		h.slaveOkNotifier.broadcast()
	}
}
