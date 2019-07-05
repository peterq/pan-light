package realtime

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

type SessionId string

type Session struct {
	Data interface{} // 业务数据存放

	id               SessionId // session 标识
	secret           string    // session 秘钥
	server           *Server
	conn             *websocket.Conn
	online           bool // 是否掉线了
	lock             sync.Mutex
	missMessage      []gson // 掉线错过的消息
	missMessageIndex int
	rooms            []*Room
}

func (ss *Session) Rooms() []*Room {
	return ss.rooms
}

func randomStr(lenght int) string {
	arr := make([]byte, lenght)
	src := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"
	for i := 0; i < lenght; i++ {
		arr[i] = byte(src[rand.Intn(len(src))])
	}
	return string(arr)
}

func newSession(conn *websocket.Conn, missMessageSize int, server *Server) *Session {
	s := &Session{
		id:          SessionId(fmt.Sprint(time.Now().UnixNano())),
		secret:      randomStr(16),
		server:      server,
		conn:        conn,
		online:      true,
		missMessage: make([]gson, missMessageSize),
	}
	return s
}

func (ss *Session) Id() SessionId {
	return ss.id
}

func (ss *Session) Emit(event string, data interface{}, room ...string) {
	d := gson{
		"type":    "event",
		"event":   event,
		"payload": data,
	}
	if len(room) > 0 {
		d["room"] = room[0]
	}
	ss.write(d)
}

func (ss *Session) write(data gson) (err error) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	if !ss.online {
		ss.missMessageIndex = (ss.missMessageIndex + 1) % len(ss.missMessage)
		ss.missMessage[ss.missMessageIndex] = gson{
			"time":    time.Now().Unix(),
			"message": data,
		}
		return errors.New("connection loss")
	}
	bin, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	if enc {
		err = websocket.Message.Send(ss.conn, encBin(bin))
	} else {
		_, err = ss.conn.Write(bin)
	}
	return
}

func (ss *Session) Read() (data gson, err error) {
	bin, err := receiveFullFrame(ss.conn)
	if err != nil {
		return
	}
	err = json.Unmarshal(bin, &data)
	return
}
func (ss *Session) InRoom(name string) bool {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	for _, room := range ss.rooms {
		if room.name == name {
			return true
		}
	}
	return false
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
