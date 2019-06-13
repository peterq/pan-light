package realtime

import (
	"fmt"
	"github.com/peterq/pan-light/server/timewheel"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

type gson = map[string]interface{}
type EventHandler interface {
	HandleEvent(ss *Session, data interface{})
}
type RpcHandler interface {
	HandleRpc(ss *Session, p gson) (result interface{}, err error)
}

type EventHandleFunc func(ss *Session, data interface{})

func (fn EventHandleFunc) HandleEvent(ss *Session, data interface{}) {
	fn(ss, data)
}

type RpcHandleFunc func(ss *Session, p gson) (result interface{}, err error)

func (fn RpcHandleFunc) HandleRpc(ss *Session, p gson) (result interface{}, err error) {
	return fn(ss, p)
}

type Server struct {
	SessionKeepTime         time.Duration // 断线, 回话维持时间
	KeepMessageCount        int           // 断线保留的消息数量
	BeforeAcceptSession     func(ss *Session) (err error)
	AfterAcceptSession      func(ss *Session) (err error)
	BeforeDispatchUserEvent func(ss *Session, event string) (err error)
	BeforeDispatchUserRpc   func(ss *Session, method string) (err error)
	OnSessionLost           func(ss *Session)

	eventHandlerMap map[string]EventHandler
	rpcHandlerMap   map[string]RpcHandler
	inited          bool

	wheel                     *timewheel.TimeWheel
	sessionMap                map[SessionId]*Session
	sessionMapLock            sync.RWMutex
	handlerMap                map[string]func(data interface{})
	deliveryTaskRemoveSession func(session *Session)
	deliveryTaskPingSession   func(session *Session)
	roomMap                   map[string]*Room
	roomMapLock               sync.RWMutex
}

type task struct {
	task    string
	payload interface{}
}
type DeliveryTaskFunc func(delay time.Duration, key interface{}, data interface{})

func (s *Server) init() {
	if s.inited {
		return
	}
	s.eventHandlerMap = map[string]EventHandler{}
	s.rpcHandlerMap = map[string]RpcHandler{}
	s.handlerMap = map[string]func(data interface{}){}
	s.sessionMap = map[SessionId]*Session{}
	s.roomMap = map[string]*Room{}
	s.inited = true
	s.wheel = timewheel.New(time.Second, 512, func(i interface{}) {
		t := i.(*task)
		hdlr := s.handlerMap[t.task]
		hdlr(t.payload)
	})
	s.wheel.Start()
	addSessionRemoveTask := s.RegisterTaskHandler("session.remove.later", func(data interface{}) {
		ss := data.(*Session)
		ss.lock.Lock()
		defer ss.lock.Unlock()
		if !ss.online {
			go s.onSessionExpired(ss)
		}
	})
	s.deliveryTaskRemoveSession = func(session *Session) {
		addSessionRemoveTask(s.SessionKeepTime, time.Now().UnixNano(), session)
	}
	s.deliveryTaskPingSession = func() func(session *Session) {
		fn := s.RegisterTaskHandler("session.ping.later", func(data interface{}) {
			ss := data.(*Session)
			ss.Emit("ping", "")
		})
		return func(session *Session) {
			fn(10*time.Second, time.Now().Unix(), session)
		}
	}()
}

func (s *Server) RegisterTaskHandler(taskName string, handler func(data interface{})) DeliveryTaskFunc {
	s.handlerMap[taskName] = handler
	return func(delay time.Duration, key interface{}, data interface{}) {
		s.wheel.AddTimer(delay, key, &task{
			task:    taskName,
			payload: data,
		})
	}
}

func (s *Server) RegisterEventHandler(mp map[string]EventHandler) {
	for key, handler := range mp {
		s.eventHandlerMap[key] = handler
	}
}

func (s *Server) RegisterRpcHandler(mp map[string]RpcHandler) {
	for key, handler := range mp {
		s.rpcHandlerMap[key] = handler
	}
}

// 用于绑定外部http server
func (s *Server) HttpHandler() http.Handler {
	s.init()
	hd := websocket.Handler(s.handleWsConn)
	return hd
}

// session 过期
func (s *Server) onSessionExpired(ss *Session) {
	s.RemoveSession(ss.id)
}

func (s *Server) RemoveSession(id SessionId) {
	s.sessionMapLock.Lock()
	defer s.sessionMapLock.Unlock()
	ss, ok := s.sessionMap[id]
	if !ok {
		return
	}
	delete(s.sessionMap, id)
	go func() {
		for _, room := range ss.rooms {
			room.Remove(id)
		}
		if ss.online {
			ss.conn.Close()
		}
	}()
	if s.OnSessionLost != nil {
		go s.OnSessionLost(ss)
	}

}

func (s *Server) SessionById(id SessionId) (ss *Session, ok bool) {
	s.sessionMapLock.RLock()
	defer s.sessionMapLock.RUnlock()
	ss, ok = s.sessionMap[id]
	return ss, ok
}

func (s *Server) RoomByName(name string) *Room {
	s.roomMapLock.Lock()
	defer s.roomMapLock.Unlock()
	room, ok := s.roomMap[name]
	if !ok {
		room = &Room{
			name:    name,
			server:  s,
			members: sessionIdSlice{},
		}
		s.roomMap[name] = room
	}
	return room
}

func (s *Server) RoomExist(name string) bool {
	s.roomMapLock.Lock()
	defer s.roomMapLock.Unlock()
	_, ok := s.roomMap[name]
	return ok
}

func (s *Server) handleWsConn(conn *websocket.Conn) {
	log.Println("new ws conn", conn.RemoteAddr())
	var err error
	defer func() {
		conn.Close()
		log.Println("close ws conn", conn.RemoteAddr().String(), err)
		if err != io.EOF {
			log.Printf("stack %s", debug.Stack())
		}

		if e := recover(); e != nil {
			log.Println("ws connection error", e)
			log.Printf("stack %s", debug.Stack())
		}
	}()
	// 新session或者恢复之前的session
	d, err := (&Session{
		conn: conn,
	}).Read()
	if err != nil {
		return
	}
	t := d["type"].(string)
	var session *Session
	var isNewSession bool
	if t == "session.new" {
		isNewSession = true
	} else if t == "session.resume" {
		func() {
			s.sessionMapLock.RLock()
			defer s.sessionMapLock.RUnlock()
			var ok bool
			var id string
			id = d["sessionId"].(string)
			session, ok = s.sessionMap[SessionId(id)]
			if !ok || session.secret != d["sessionSecret"] {
				isNewSession = true
			} else {
				session.conn.Close()
				session.conn = conn
				session.online = true
			}
		}()
	} else {
		err = errors.New("session handshake error")
		return
	}

	if isNewSession {
		session = newSession(conn, s.KeepMessageCount, s)
		if s.BeforeAcceptSession != nil {
			if err = s.BeforeAcceptSession(session); err != nil {
				return
			}
		}
		s.sessionMapLock.Lock()
		s.sessionMap[session.id] = session
		s.sessionMapLock.Unlock()
		session.Emit("session.new", gson{
			"id":     session.id,
			"secret": session.secret,
		})
		if s.AfterAcceptSession != nil {
			if err = s.AfterAcceptSession(session); err != nil {
				log.Println(err)
			}
		}
	} else {
		session.Emit("session.resume", "ok")
	}

	defer func() {
		session.lock.Lock()
		defer session.lock.Unlock()
		session.online = false
		s.deliveryTaskRemoveSession(session)
	}()
	err = s.readMessageLoop(session)
}

func (s *Server) readMessageLoop(ss *Session) (err error) {
	var data gson
	for {
		data, err = ss.Read()
		if err != nil {
			if err != io.EOF {
				err = errors.Wrap(err, "read error")
			}
			return
		}
		go s.handleMessage(data, ss)
	}
}

func (s *Server) handleMessage(data gson, ss *Session) {
	defer func() {
		if e := recover(); e != nil {
			log.Println("handle message error", e, data)
			log.Printf("stack %s", debug.Stack())
			ss.Emit("server.handle.error", gson{
				"sourceMessage": data,
				"error":         fmt.Sprint(e),
			})
		}
	}()

	if data["type"] == "call" {
		defer func() {
			if e := recover(); e != nil {
				log.Println("rpc error", e, data)
				log.Printf("stack %s", debug.Stack())
				ss.write(gson{
					"type":    "call.result",
					"success": false,
					"error":   fmt.Sprint(e),
					"id":      data["id"],
				})
			}
		}()
		method := data["method"].(string)
		if s.BeforeDispatchUserRpc != nil && s.BeforeDispatchUserRpc(ss, method) != nil {
			return
		}
		handler, ok := s.rpcHandlerMap[method]
		if !ok {
			ss.write(gson{
				"type":    "call.result",
				"success": false,
				"error":   "rpc error, handler not defined",
				"id":      data["id"],
			})
			return
		}
		result, err := handler.HandleRpc(ss, data["param"].(gson))
		resp := gson{
			"type":    "call.result",
			"success": err == nil,
			"result":  result,
			"id":      data["id"],
		}
		if err != nil {
			resp["error"] = err.Error()
		}
		ss.write(resp)
		return
	}
	if data["type"] == "event" {
		event := data["event"].(string)
		if s.BeforeDispatchUserEvent != nil && s.BeforeDispatchUserEvent(ss, event) != nil {
			return
		}
		handler := s.eventHandlerMap[event]
		handler.HandleEvent(ss, data["payload"])
		return
	}
}
