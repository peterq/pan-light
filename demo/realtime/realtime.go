package realtime

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

type gson = map[string]interface{}

type callResult struct {
	ret interface{}
	err error
}

type RealTime struct {
	WsAddr       string
	Role         string
	HostName     string
	HostPassWord string
	SlaveName    string // role 为 slave 时可用
	OnConnected  func()

	inited        bool
	conn          *websocket.Conn
	connectOK     bool
	connectLock   sync.Mutex
	connectOkCond *sync.Cond
	sessionId     string
	sessionSecret string

	listenerMap map[string][]func(data interface{}, room string)
	callMap     sync.Map
	//callMap     map[float64]chan<- *callResult
	callMapLock sync.Mutex
	logWsMsg    bool
}

func (rt *RealTime) Init() {
	if rt.inited {
		return
	}
	rt.inited = true
	rt.connectOkCond = sync.NewCond(&rt.connectLock)
	rt.listenerMap = map[string][]func(data interface{}, room string){}
	go func() {
		for range time.Tick(10 * time.Second) {
			rt.Call("ping", gson{})
		}
	}()
	rt.RegisterEventListener(map[string]func(data interface{}, room string){
		"session.new": func(data interface{}, room string) {
			rt.sessionId = data.(gson)["id"].(string)
			rt.sessionSecret = data.(gson)["id"].(string)
			log.Println("rt 会话创建成功")
		},
		"session.resume": func(data interface{}, room string) {
			log.Println("rt 会话恢复成功")
		},
	})
	go rt.connect()
}

func (rt *RealTime) RegisterEventListener(mp map[string]func(data interface{}, room string)) {
	for event, fn := range mp {
		if _, ok := rt.listenerMap[event]; !ok {
			rt.listenerMap[event] = []func(interface{}, string){fn}
		} else {
			rt.listenerMap[event] = append(rt.listenerMap[event], fn)
		}
	}
}

func (rt *RealTime) connect() {
	go func() {
		for {
			rt.connectLock.Lock()
			for rt.connectOK != true { // 等待连接ok
				rt.connectOkCond.Wait()
			}
			rt.connectLock.Unlock()
			if rt.OnConnected != nil {
				go rt.OnConnected()
			}
			log.Println("ws连接成功")
			rt.readLoop()
		}
	}()
	first := true
	for {
		func() { // 等待链接不ok, 进行连接
			rt.connectLock.Lock()
			for rt.connectOK != false {
				rt.connectOkCond.Wait()
			}
			defer func() {
				rt.connectLock.Unlock()
				rt.connectOkCond.Broadcast()
			}()
			if !first {
				time.Sleep(5 * time.Second)
			}
			first = false
			var err error
			conn, err := websocket.Dial(rt.WsAddr, "", "http://localhost/")
			if err != nil {
				log.Println("ws 连接错误", err)
				time.Sleep(5 * time.Second)
				return
			}
			rt.conn = conn
			rt.connectOK = true
			if rt.sessionId != "" {
				log.Println("will resume session")
				err = rt.write(gson{
					"type":          "session.resume",
					"sessionId":     rt.sessionId,
					"sessionSecret": rt.sessionSecret,
				})
			} else {
				log.Println("will new session")
				err = rt.write(gson{
					"type": "session.new",
				})
				err = rt.write(gson{
					"role":        rt.Role,
					"host_name":   rt.HostName,
					"host_secret": rt.HostPassWord,
					"slave_name":  rt.SlaveName,
				})
			}
			if err != nil {
				log.Println("write error", err)
				rt.connectOK = false
				return
			}
		}()
	}
}

func (rt *RealTime) readLoop() {
	for {
		data, err := rt.read()
		if err != nil {
			log.Println("read ws error, connect again", err)
			prevCon := rt.conn
			rt.connectLock.Lock()
			if prevCon == rt.conn {
				rt.connectOK = false
				rt.connectOkCond.Broadcast()
			}
			rt.connectLock.Unlock()
			return
		}
		go rt.handleMsg(data)
	}
}

func (rt *RealTime) handleMsg(data gson) {
	t := data["type"].(string)
	if t == "event" {
		event := data["event"].(string)
		room, ok := data["room"]
		if !ok {
			room = ""
		}
		cbs, ok := rt.listenerMap[event]
		if !ok {
			return
		}
		for _, cb := range cbs {
			go func(cb func(data interface{}, room string)) {
				defer func() {
					if e := recover(); e != nil {
						log.Println(e)
						debug.PrintStack()
					}
				}()
				cb(data["payload"], room.(string))
			}(cb)
		}
		return
	}
	if t == "call.result" {
		id := data["id"].(float64)
		ch, ok := rt.callMap.Load(id)
		if !ok {
			return
		}
		rt.callMapLock.Lock()
		defer rt.callMapLock.Unlock()
		ret := &callResult{}
		if data["success"].(bool) {
			ret.err = nil
			ret.ret = data["result"]
		} else {
			ret.err = errors.New(data["error"].(string))
			ret.ret = nil
		}
		ch.(chan *callResult) <- ret
		return
	}

}

func (rt *RealTime) Emit(event string, data interface{}) {
	rt.write(gson{
		"type":    "event",
		"event":   rt.Role + "." + event,
		"payload": data,
	})
}

func (rt *RealTime) Broadcast(room string, event string, data interface{}) {
	rt.Emit("broadcast", gson{
		"room":    room,
		"event":   event,
		"payload": data,
	})
}

func (rt *RealTime) Call(method string, param gson) (result interface{}, err error) {

	id := float64(time.Now().UnixNano())
	ch := make(chan *callResult)

	rt.callMapLock.Lock()
	rt.callMap.Store(id, ch)
	rt.callMapLock.Unlock()

	rt.write(gson{
		"type":   "call",
		"method": rt.Role + "." + method,
		"param":  param,
		"id":     id,
	})

	ret := <-ch
	result, err = ret.ret, ret.err
	close(ch)

	rt.callMapLock.Lock()
	rt.callMap.Delete(id)
	rt.callMapLock.Unlock()

	return
}

func (rt *RealTime) read() (data gson, err error) {
	bin, err := receiveFullFrame(rt.conn)
	if err != nil {
		return
	}
	err = json.Unmarshal(bin, &data)
	if rt.logWsMsg {
		log.Println("ws <-", data)
	}
	return
}

func (rt *RealTime) write(data gson) (err error) {

	if !rt.connectOK {
		return errors.New("connect not ok")
	}

	if rt.logWsMsg {
		log.Println("ws ->", data)
	}
	bin, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	if enc {
		err = websocket.Message.Send(rt.conn, encBin(bin))
	} else {
		_, err = rt.conn.Write(bin)
	}
	if err != nil {
		log.Println("write ws error, connect again", err)
		prevCon := rt.conn
		rt.connectLock.Lock()
		if prevCon == rt.conn {
			rt.connectOK = false
			rt.connectOkCond.Broadcast()
		}
		rt.connectLock.Unlock()
		return
	}
	return
}

const enc = true
const key = "pan-light"

func xorBin(bin []byte) []byte {
	dst := make([]byte, len(bin))
	keyLen := len(key)
	for idx, b := range bin {
		dst[idx] = key[idx%keyLen] ^ b
	}
	return dst
}

func encBin(bin []byte) []byte {
	if !enc {
		return bin
	}
	zipped, _ := gzipEncode(bin)
	return xorBin(zipped)
}
func dencBin(bin []byte) (dest []byte, err error) {
	if !enc {
		return bin, nil
	}
	return gzipDecode(xorBin(bin))
}

func gzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

func gzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

// 接受完整帧
func receiveFullFrame(ws *websocket.Conn) ([]byte, error) {
	var data []byte
	for {
		var seg []byte
		fin, err := receiveFrame(websocket.Message, ws, &seg)
		if err != nil {
			return nil, err
		}
		data = append(data, seg...)
		if fin {
			break
		}
	}
	return dencBin(data)
}

// 接受帧
func receiveFrame(cd websocket.Codec, ws *websocket.Conn, v interface{}) (fin bool, err error) {
again:
	frame, err := ws.NewFrameReader()
	if frame.HeaderReader() != nil {
		bin := make([]byte, 1)
		frame.HeaderReader().Read(bin)
		fin = ((bin[0] >> 7) & 1) != 0
	}
	if err != nil {
		return
	}
	frame, err = ws.HandleFrame(frame)
	if err != nil {
		return
	}
	if frame == nil {
		goto again
	}

	payloadType := frame.PayloadType()

	data, err := ioutil.ReadAll(frame)
	if err != nil {
		return
	}
	return fin, cd.Unmarshal(data, payloadType, v)
}
