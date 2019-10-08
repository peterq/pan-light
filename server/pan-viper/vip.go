package pan_viper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/peterq/pan-light/server/artisan"
	"github.com/peterq/pan-light/server/dao"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var vipMap sync.Map

func init() {
	data, err := dao.Vip.GetAll()
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		panic("no vip account available")
	}
	for _, model := range data {
		v := &Vip{
			username:  model.Username,
			cookieRaw: model.Cookie,
		}
		v.init()
		vipMap.Store(model.Username, v)
	}
}

type loginSession struct {
	Sign      string
	Timestamp string
	Bdstoken  string
	Bduss     string
	createAt  time.Time
}

type Vip struct {
	http             http.Client
	username         string
	bduss            string
	cookieRaw        string
	_loginSession    *loginSession
	loginSessionLock sync.RWMutex
}

func (v *Vip) loginSession() *loginSession {
	v.loginSessionLock.RLock()
	defer v.loginSessionLock.RUnlock()
	return v._loginSession
}

func (v *Vip) init() {
	v.http, v.bduss = makeHttpClient(v.cookieRaw)
	go v.CreateSession()
}

func (v *Vip) Username() string {
	return v.username
}

func (v *Vip) CreateSession() (err error) {
	defer func() {
		if err != nil {
			artisan.App.Logger().Error(v.username, "vip创建session错误", err)
			go v.CreateSession()
		}
	}()
	old := v.loginSession()
	if old != nil && time.Now().Sub(old.createAt) < time.Second {
		return
	}
	v.loginSessionLock.Lock()
	defer v.loginSessionLock.Unlock()
	// 高并发下防止重复更新session
	if old != v._loginSession {
		return
	}
	homePageLink := "https://pan.baidu.com/disk/homePageLink"
	req := newRequest("GET", homePageLink)
	res, err := v.http.Do(req)
	if err != nil {
		err = errors.Wrap(err, "访问首页错误")
		return
	}
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "读取首页body错误")
		return
	}
	body := string(bin)
	reg := regexp.MustCompile(`var context=(.*);\n`)
	find := reg.FindStringSubmatch(body)
	if res.Request.URL.String() != homePageLink {
		err = errors.New("重定向到" + res.Request.URL.String())
		return
	}
	raw := gson{}
	err = json.Unmarshal([]byte(find[1]), &raw)
	if err != nil {
		log.Println(body)
		return
	}
	s := loginSession{
		Sign:      "",
		Timestamp: "",
		Bdstoken:  "",
		Bduss:     "",
	}
	s.Sign = loginSign(raw["sign3"].(string), raw["sign1"].(string))
	s.Timestamp = fmt.Sprint(int(raw["timestamp"].(float64)))
	s.Bdstoken = raw["bdstoken"].(string)
	s.Bduss = v.bduss
	s.createAt = time.Now()
	v._loginSession = &s
	log.Println(v.username, "完成loginSession")
	return
}

func loginSign(j, r string) string {
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

func (v *Vip) DeleteFile(serverPath string) (err error) {
	ss := v.loginSession()
	data, err := v.request("POST", "https://pan.baidu.com/api/filemanager", gson{
		"opera":      "delete",
		"async":      2,
		"onnest":     "fail",
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   ss.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"filelist": fmt.Sprintf("[\"%s\"]", serverPath),
	})
	if err != nil {
		return
	}
	data, err = v.request("POST", "https://pan.baidu.com/share/taskquery", gson{
		"taskid":     int64(data["taskid"].(float64)),
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"bdstoken":   ss.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": "0",
	}, gson{
		"filelist": fmt.Sprintf("[\"%s\"]", serverPath),
	})
	return
}

func (v *Vip) SaveFileByMd5(md5, sliceMd5, path string, contentLength int64) (fid string, fileSize int64, err error) {
	err = errors.New("vip通道已关闭, 请联系管理员")
	return
	ss := v.loginSession()
	data, err := v.request("POST", "https://pan.baidu.com/api/rapidupload", gson{
		"rtype":      1,
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   ss.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"path":           path,
		"content-length": contentLength,
		"content-md5":    md5,
		"slice-md5":      sliceMd5,
		"target_path":    filepath.Dir(path),
		"local_mtime":    1533345687,
	})
	if err != nil {
		err = errors.Wrap(err, "极速上传到vip账号失败")
		return
	}
	if _, ok := data["errno"]; !ok {
		log.Println(data)
		err = errors.New("极速上传到vip账号失败")
	}
	info := data["info"].(gson)
	fid = fmt.Sprint(int64(info["fs_id"].(float64)))
	fileSize = int64(info["size"].(float64))
	serverPath := info["path"].(string)
	if serverPath[len(serverPath)-1] == ')' {
		go v.DeleteFile(path)
	}
	return
}

func (v *Vip) LinkByFid(fid string) (link string, err error) {
	ss := v.loginSession()

	data, err := v.request("GET", "https://pan.baidu.com/api/download", gson{
		"sign":       ss.Sign,
		"timestamp":  ss.Timestamp,
		"fidlist":    "[" + fid + "]",
		"type":       "dlink",
		"channel":    "chunlei",
		"web":        5,
		"app_id":     "250528",
		"bdstoken":   ss.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 5,
	}, nil)
	if err != nil {
		err = errors.Wrap(err, "")
		return
	}
	link = data["dlink"].([]interface{})[0].(map[string]interface{})["dlink"].(string)
	link = v.getRedirectedLink(link)
	return
}

func (v *Vip) getRedirectedLink(link string) string {
	req := newRequest("GET", link)
	resp, err := v.http.Do(req)
	if err != nil {
		log.Println(err)
	}
	end := resp.Request.URL.String()
	resp.Body.Close()
	return end
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
		// 页面过期错误码
		if n.(float64) == 112 {
			go v.CreateSession()
		}
		err = errors.New("pan api error code " + fmt.Sprint(data["errno"]))
	}
	return
}

func GetVip() *Vip {
	var v *Vip
	vipMap.Range(func(key, value interface{}) bool {
		v = value.(*Vip)
		return false
	})
	return v
}

func GetVipByUsername(username string) *Vip {
	if v, ok := vipMap.Load(username); ok {
		return v.(*Vip)
	}
	return nil
}
