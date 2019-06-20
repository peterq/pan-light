package pan_viper

import (
	"encoding/json"
	"fmt"
	"github.com/peterq/pan-light/server/dao"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var vipMap = map[string]*Vip{}

func init() {
	data, err := dao.Vip.GetAll()
	if err != nil {
		panic(err)
	}
	for _, model := range data {
		vipMap[model.Username] = &Vip{
			http:     makeHttpClient(model.Bduss),
			username: model.Username,
		}
	}
}

type loginSession struct {
}

type Vip struct {
	http             http.Client
	username         string
	_loginSession    loginSession
	loginSessionLock sync.RWMutex
}

func (v *Vip) loginSession() loginSession {
	v.loginSessionLock.RLock()
	defer v.loginSessionLock.RUnlock()
	return v._loginSession
}

func (v *Vip) Username() string {
	return v.username
}

func (v *Vip) CreateSession() {
	//v.request()
}

func (v *Vip) LoadShareFilenameAndUk(link, secret string) (uk, filename string, share gson, err error) {
	err = v.inputSharePwd(link, secret)
	if err != nil {
		err = errors.Wrap(err, "input pwd error")
	}
	resp, err := v.http.Get(link)
	if err != nil {
		err = errors.Wrap(err, "load share page error")
		return
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "load share page error")
		return
	}
	body := string(bin)
	reg := regexp.MustCompile(`setData\(({.*})\);\n`)
	find := reg.FindStringSubmatch(body)
	if len(find) != 2 {
		err = errors.New("load share info error")
	}
	err = json.Unmarshal([]byte(find[1]), &share)
	if err != nil {
		err = errors.Wrap(err, "decode share info error")
		return
	}
	uk = fmt.Sprint(int64(share["uk"].(float64)))
	code := int64(share["file_list"].(gson)["errno"].(float64))
	if code != 0 {
		err = errors.New("share error code: " + fmt.Sprint(code))
		return
	}
	filename = share["file_list"].(gson)["list"].([]interface{})[0].(gson)["server_filename"].(string)
	return
}

func (v *Vip) inputSharePwd(link, secret string) (err error) {

	t := strings.Split(link, "/")
	surl := t[len(t)-1]
	t = strings.Split(surl, "")

	t = t[1:]
	surl = strings.Join(t, "")

	_, err = v.request("POST", "https://pan.baidu.com/share/verify", gson{
		"surl":       surl,
		"t":          time.Now().UnixNano() / int64(time.Millisecond),
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   "null",
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"pwd":       secret,
		"vcode":     "",
		"vcode_str": "",
	})
	return
}

func (v *Vip) request(method, link string, params gson, form gson) (data gson, err error) {
	req := newRequest(method, link, form)
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
	}
	req.URL.RawQuery = q.Encode()
	resp, err := v.http.Do(req)

	if err != nil {
		err = errors.Wrap(err, "http request error")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "read http resp error")
		return
	}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return
	}
	if n, ok := data["errno"]; ok && n.(float64) != 0 {
		err = errors.New("pan api error code " + fmt.Sprint(data["errno"]))
	}
	return
}

func GetVip() *Vip {
	for _, value := range vipMap {
		return value
	}
	return nil
}
