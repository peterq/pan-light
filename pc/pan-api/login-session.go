package pan_api

import (
	"encoding/base64"
	"fmt"
	"github.com/peterq/pan-light/pc/storage"
)

type LoginSessionStruct struct {
	rawJson *tJson

	Avatar    string
	Username  string
	Sign      string
	Timestamp string
	Bdstoken  string
	Bduss     string
}

var LoginSession *LoginSessionStruct
var bduss string

func handleLoginSession(rawJson *tJson) {
	LoginSession = new(LoginSessionStruct)
	LoginSession.rawJson = rawJson
	raw := *rawJson
	LoginSession.Avatar = raw["photo"].(string)
	LoginSession.Username = raw["username"].(string)
	LoginSession.Sign = sign(raw["sign3"].(string), raw["sign1"].(string))
	LoginSession.Timestamp = fmt.Sprint(int(raw["timestamp"].(float64)))
	LoginSession.Bdstoken = raw["bdstoken"].(string)
	LoginSession.Bduss = bduss
	storage.UserState.Uk = fmt.Sprint(int64(raw["uk"].(float64)))
}
func sign(j, r string) string {
	a := [256]int{}
	p := [256]int{}
	o := make([]byte, len(r))
	v := len(j)
	for q := 0; q < 256; q++ {
		a[q] = int(j[q%v : q%v+1][0])
		p[q] = q
	}
	for u, q := 0, 0; q < 256; q++ {
		u = (u + p[q] + a[q]) % 256
		t := p[q]
		p[q] = p[u]
		p[u] = t
	}
	for i, u, q := 0, 0, 0; q < len(r); q++ {
		i = (i + 1) % 256
		u = (u + p[i]) % 256
		t := p[i]
		p[i] = p[u]
		p[u] = t
		k := p[((p[i] + p[u]) % 256)]
		o[q] = byte(int(r[q : q+1][0]) ^ k)
	}
	return base64.StdEncoding.EncodeToString(o)
}
