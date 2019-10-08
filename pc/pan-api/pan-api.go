package pan_api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/peterq/pan-light/pc/storage"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var urlMap = map[string]string{
	"home":  "https://pan.baidu.com/disk/home",
	"list":  "https://pan.baidu.com/api/list",
	"dlink": "https://pan.baidu.com/api/download",
	"usage": "https://pan.baidu.com/api/quota",
}

func newRequest(method, link string, body ...gson) *http.Request {
	if u, ok := urlMap[link]; ok {
		link = u
	}
	var bd io.Reader
	if len(body) == 1 {
		formData := url.Values{}
		for key, value := range body[0] {
			formData.Add(key, fmt.Sprint(value))
		}
		bd = strings.NewReader(formData.Encode())
	}
	req, err := http.NewRequest(method, link, bd)
	req.Header.Set("user-agent", BaiduUA)
	if err != nil {
		log.Println(err)
	}
	return req
}

func GetSign() (ctx map[string]interface{}, err error) {
	if storage.UserState.PanCookie != nil {
		cookieJar, _ = cookiejar.New(nil)
		httpClient.Jar = cookieJar

		var cookies []*http.Cookie
		for _, c := range storage.UserState.PanCookie {
			cookies = append(cookies, &http.Cookie{
				Name:   c.Key,
				Value:  c.Value,
				Domain: ".baidu.com",
			})
			if c.Key == "BDUSS" {
				bduss = c.Value
			}
		}
		u, _ := url.Parse("https://pan.baidu.com")
		cookieJar.SetCookies(u, cookies)
	}
	req := newRequest("GET", "home")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(req.Cookies())
		log.Println(err)
		return
	}
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(string(bin))
		return
	}
	body := string(bin)
	reg := regexp.MustCompile(`var context=(.*);\n`)
	find := reg.FindStringSubmatch(body)
	if res.Request.URL.String() != urlMap["home"] {
		log.Println("重定向: ", res.Request.URL.String())
		err = errors.New("未登录")
		return
	}
	raw := tJson{}
	err = json.Unmarshal([]byte(find[1]), &raw)
	if err != nil {
		log.Println(body)
		return
	}
	handleLoginSession(&raw)
	ctx = raw
	return
}

// 获取网盘使用空间
func Usage() (result interface{}, err error) {
	req := newRequest("GET", "usage")
	params := map[string]interface{}{
		"checkexpire": 1,
		"checkfree":   1,
		"channel":     "chunlei",
		"web":         1,
		"app_id":      250528,
		"bdstoken":    LoginSession.Bdstoken,
		"logid":       time.Now().UnixNano(),
		"clienttype":  0,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取磁盘空间错误, 错误码" + fmt.Sprint(data["errno"]))
	}
	result = data
	return
}

// 获取目录下的文件(夹)
func ListDir(path string) (list interface{}, err error) {
	req := newRequest("GET", "list")
	params := map[string]interface{}{
		"channel":    "chunlei",
		"clienttype": 0,
		"web":        1,
		"showempty":  1,
		"num":        10000,
		"t":          time.Now().Unix() * 1000,
		"dir":        path,
		"page":       1,
		"desc":       1,
		"order":      "name",
		"_":          time.Now().Unix() * 1000,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取文件夹信息错误, 错误码" + fmt.Sprint(data["errno"]))
	}
	list = data["list"]
	return
}

// 链接解析
func LinkDirect(fid string) (link string, err error) {

	req := newRequest("GET", "dlink")
	params := map[string]interface{}{
		"sign":       LoginSession.Sign,
		"timestamp":  LoginSession.Timestamp,
		"fidlist":    "[" + fid + "]",
		"type":       "dlink",
		"channel":    "chunlei",
		"web":        5,
		"app_id":     "250528",
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 5,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	log.Println(req.URL.String())
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	var data tJson
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if data["errno"].(float64) != 0 {
		err = errors.New("获取文件信息错误, 错误码" + fmt.Sprint(data["errno"]))
		return
	}
	link = data["dlink"].([]interface{})[0].(map[string]interface{})["dlink"].(string)
	link = getRedirectedLink(link)
	return
}

func getRedirectedLink(link string) string {
	req := newRequest("GET", link)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	end := resp.Request.URL.String()
	log.Println(end)
	resp.Body.Close()
	return end
}

func md5bin(bin []byte) string {
	h := md5.New()
	h.Write(bin)
	return hex.EncodeToString(h.Sum(nil))
}

func UploadText(content string, path string) (fid, serverPath string, err error) {
	bin := []byte(content)
	md5str := md5bin(bin)
	// pre create
	data, err := request("POST", "https://pan.baidu.com/api/precreate", gson{
		"channel":      "chunlei",
		"web":          1,
		"app_id":       250528,
		"bdstoken":     LoginSession.Bdstoken,
		"logid":        time.Now().UnixNano(),
		"clienttype":   0,
		"startLogTime": time.Now().UnixNano() / int64(time.Millisecond),
	}, gson{
		"path":        path,
		"autoinit":    1,
		"target_path": "/",
		"block_list":  fmt.Sprintf("[\"%s\"]", md5str),
		"local_mtime": time.Now().Unix(),
	})
	if err != nil {
		return
	}
	uploadId := data["uploadid"].(string)
	// upload
	bodyBuf := bytes.NewBuffer([]byte{})
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "1.txt")
	if err != nil {
		return
	}
	fileWriter.Write(bin)
	bodyWriter.Close()

	req, err := http.NewRequest("POST", "https://qdcu01.pcs.baidu.com/rest/2.0/pcs/superfile2", bodyBuf)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Set("user-agent", BaiduUA)

	params := map[string]interface{}{
		"method":     "upload",
		"app_id":     250528,
		"channel":    "chunlei",
		"clienttype": 0,
		"web":        1,
		"BDUSS":      LoginSession.Bduss,
		"logid":      time.Now().UnixNano(),
		"path":       path,
		"uploadid":   uploadId,
		"uploadsign": 0,
		"partseq":    0,
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	data, err = sendRequest(req)
	if err != nil {
		return
	}
	if _, ok := data["md5"]; !ok {
		err = errors.New("upload fail")
		return
	}
	blockMd5 := data["md5"].(string)
	// combine
	data, err = request("POST", "https://pan.baidu.com/api/create", gson{
		"isdir":      0,
		"rtype":      1,
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"path":        path,
		"size":        len(bin),
		"uploadid":    uploadId,
		"autoinit":    1,
		"target_path": "/",
		"block_list":  fmt.Sprintf("[\"%s\"]", blockMd5),
		"local_mtime": time.Now().Unix(),
	})
	if err != nil {
		return
	}
	fid = fmt.Sprint(int64(data["fs_id"].(float64)))
	serverPath = data["path"].(string)
	return
}

type gson = map[string]interface{}

func request(method, link string, params gson, form gson) (data gson, err error) {
	req := newRequest(method, link, form)
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
	}
	req.URL.RawQuery = q.Encode()
	return sendRequest(req)
}

func sendRequest(req *http.Request) (data gson, err error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	body := readHtml(resp.Body)
	err = json.Unmarshal(tBin(body), &data)
	if err != nil {
		return
	}
	if n, ok := data["errno"]; ok && n.(float64) != 0 {
		err = errors.New("错误码" + fmt.Sprint(data["errno"]))
	}
	return
}

func ShareFile(fid, secret string) (link, sec string, err error) {
	if len(secret) != 4 {
		secret = randomStr(4)
	}
	sec = secret
	data, err := request("POST", "https://pan.baidu.com/share/set", gson{
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"schannel":     4,
		"channel_list": "[]",
		"period":       7,
		"pwd":          secret,
		"fid_list":     fmt.Sprintf("[%s]", fid),
	})
	if err != nil {
		return
	}
	link = data["link"].(string)
	return
}

func randomStr(length int) string {
	arr := make([]byte, length)
	src := "qwertyuiopasdfghjklzxcvbnm1234567890"
	for i := 0; i < length; i++ {
		arr[i] = byte(src[rand.Intn(len(src))])
	}
	return string(arr)
}

func DeleteFile(serverPath string) (err error) {
	data, err := request("POST", "https://pan.baidu.com/api/filemanager", gson{
		"opera":      "delete",
		"async":      2,
		"onnest":     "fail",
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": 0,
	}, gson{
		"filelist": fmt.Sprintf("[\"%s\"]", serverPath),
	})
	if err != nil {
		return
	}
	data, err = request("POST", "https://pan.baidu.com/share/taskquery", gson{
		"taskid":     int64(data["taskid"].(float64)),
		"channel":    "chunlei",
		"web":        "1",
		"app_id":     "250528",
		"bdstoken":   LoginSession.Bdstoken,
		"logid":      time.Now().UnixNano(),
		"clienttype": "0",
	}, gson{
		"filelist": fmt.Sprintf("[\"%s\"]", serverPath),
	})
	log.Println(data, err)
	return
}

// md5转存
func SaveFileByMd5(md5, sliceMd5, path string, contentLength int64) (serverPath, fid string, fileSize int64, err error) {
	data, err := request("POST", "https://pan.baidu.com/api/rapidupload", gson{
		"rtype":      1,
		"channel":    "chunlei",
		"web":        1,
		"app_id":     250528,
		"bdstoken":   LoginSession.Bdstoken,
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
		err = errors.Wrap(err, "极速上传失败")
		return
	}
	if _, ok := data["errno"]; !ok {
		log.Println(data)
		err = errors.New("极速上传失败")
		return
	}
	info := data["info"].(gson)
	fid = fmt.Sprint(int64(info["fs_id"].(float64)))
	fileSize = int64(info["size"].(float64))
	serverPath = info["path"].(string)
	return
}
