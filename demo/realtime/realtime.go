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
	OnConnected  func()

	inited        bool
	conn          *websocket.Conn
	connectOK     bool
	connectLock   sync.Mutex
	connectOkCond *sync.Cond
	sessionId     string
	sessionSecret string

	listenerMap map[string][]func(data interface{})
	callMap     map[float64]chan<- *callResult
	callMapLock sync.Mutex
	logWsMsg    bool
}

func (rt *RealTime) Init() {
	if rt.inited {
		return
	}
	rt.inited = true
	rt.connectOkCond = sync.NewCond(&rt.connectLock)
	rt.listenerMap = map[string][]func(data interface{}){}
	rt.callMap = map[float64]chan<- *callResult{}
	go rt.connect()
}

func (rt *RealTime) RegisterEventListener(mp map[string]func(data interface{})) {
	for event, fn := range mp {
		if _, ok := rt.listenerMap[event]; !ok {
			rt.listenerMap[event] = []func(interface{}){fn}
		} else {
			rt.listenerMap[event] = append(rt.listenerMap[event], fn)
		}
	}
}

func (rt *RealTime) connect() {
	go func() {
		for {
			rt.connectLock.Lock()
			for rt.connectOK != true {
				rt.connectOkCond.Wait()
			}
			rt.connectLock.Unlock()
			if rt.OnConnected != nil {
				go rt.OnConnected()
			}
			rt.readLoop()
		}
	}()
	first := true
	for {
		func() {
			rt.connectLock.Lock()
			defer func() {
				rt.connectLock.Unlock()
				rt.connectOkCond.Broadcast()
			}()
			for rt.connectOK != false {
				rt.connectOkCond.Wait()
			}
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
				err = rt.write(gson{
					"type":          "session.resume",
					"sessionId":     rt.sessionId,
					"sessionSecret": rt.sessionSecret,
				})
			} else {
				err = rt.write(gson{
					"type": "session.new",
				})
			}
			err = rt.write(gson{
				"role":        rt.Role,
				"host_name":   rt.HostName,
				"host_secret": rt.HostPassWord,
			})
			if err != nil {
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
		cbs, ok := rt.listenerMap[event]
		if !ok {
			return
		}
		for _, cb := range cbs {
			go func() {
				defer func() {
					if e := recover(); e != nil {
						log.Println(e)
						debug.PrintStack()
					}
				}()
				cb(data["payload"])
			}()
		}
		return
	}
	if t == "call.result" {
		id := data["id"].(float64)
		ch, ok := rt.callMap[id]
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
		ch <- ret
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

func (rt *RealTime) Call(method string, param gson) (result interface{}, err error) {

	id := float64(time.Now().UnixNano())
	ch := make(chan *callResult)

	rt.callMapLock.Lock()
	rt.callMap[id] = ch
	rt.callMapLock.Unlock()

	rt.write(gson{
		"type":   "call",
		"method": method,
		"param":  param,
		"id":     id,
	})

	ret := <-ch
	result, err = ret.ret, ret.err
	close(ch)

	rt.callMapLock.Lock()
	delete(rt.callMap, id)
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
