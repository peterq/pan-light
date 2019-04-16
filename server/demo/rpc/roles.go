package rpc

import (
	"github.com/peterq/pan-light/server/realtime"
	"sync"
)

type roleType interface {
	roleName() string
}

type roleHost struct {
	name    string
	session *realtime.Session
	slaves  map[string]*roleSlave
}

func (*roleHost) roleName() string {
	return "host"
}

type roleUser struct {
	session *realtime.Session

	waitState *waitState
}

func (*roleUser) roleName() string {
	return "user"
}

type roleSlave struct {
	name        string
	host        *roleHost
	session     *realtime.Session
	userSession *realtime.Session
	lock        sync.Mutex
}

func (*roleSlave) roleName() string {
	return "slave"
}
