package nickname

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

type taskItem struct {
	book     string
	nickname string
	images   []string
	img      []byte
}

// 搜索图像
func FetchAndSaveAvatarFromInternet() {
	avatarInit()
	nicknameChan := make(chan *taskItem)
	imageChan := make(chan *taskItem)
	resultChan := make(chan *taskItem)
	searchGroup := new(sync.WaitGroup)
	searchGroup.Add(10)
	for i := 0; i < 10; i++ {
		go searchImageLoop(nicknameChan, imageChan, searchGroup)
	}
	faceCheckGroup := new(sync.WaitGroup)
	faceCheckGroup.Add(10)
	for i := 0; i < 10; i++ {
		go faceCheckLoop(imageChan, resultChan, faceCheckGroup)
	}
	nicknameChan <- &taskItem{
		book:     "飞狐外传",
		nickname: "马春花",
		images:   []string{},
	}
	close(nicknameChan)
	go func() {
		// 处理完成关闭相关通道
		searchGroup.Wait()
		close(imageChan)
		faceCheckGroup.Wait()
		close(resultChan)
	}()
	for itemResult := range resultChan {
		handleItemResult(itemResult)
	}
	log.Println("处理完成")
}

func handleItemResult(item *taskItem) {
	log.Println("result", item.nickname)
}

func searchImageLoop(nicknameChan chan *taskItem, imageChan chan *taskItem, wg *sync.WaitGroup) {
	go func() {
		for item := range nicknameChan {
			log.Println("search item", item.nickname)
			item.images = searchImage(item.book + " " + item.nickname)
			imageChan <- item
		}
		wg.Done()
	}()
	return
}

func faceCheckLoop(imageChan chan *taskItem, resultChan chan *taskItem, wg *sync.WaitGroup) {
	go func() {
		for item := range imageChan {
			log.Println("check face", item.nickname, item.images)
			for _, img := range item.images {
				resp, err := http.Get(img)
				if err != nil {
					continue
				}
				bin, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					continue
				}
				item.img = bin
				checkFace(bin)
				break
			}
			resultChan <- item
		}
		wg.Done()
	}()
	return
}

var searchHttpClient http.Client
var faceHttpClient http.Client

func avatarInit() {
	jar, _ := cookiejar.New(nil)
	jar1, _ := cookiejar.New(nil)
	parallel := 20
	searchHttpClient = http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    parallel,
			MaxConnsPerHost: parallel,
		},
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       0,
	}
	faceHttpClient = searchHttpClient
	faceHttpClient.Jar = jar1
	req := &http.Request{}
	req.URL, _ = url.Parse("http://kan.msxiaobing.com/V3/Portal?task=yanzhi&ftid=")
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("Referer", "http://kan.msxiaobing.com/V3/Portal?task=yanzhi&ftid=")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Mobile Safari/537.36")
	faceHttpClient.Do(req)
}

type gson = map[string]interface{}

func searchImage(keyword string) (result []string) {
	req := &http.Request{}
	req.URL, _ = url.Parse("https://m.baidu.com/sf/vsearch/image/search/wisesearchresult")
	params := map[string]interface{}{
		"tn":         "wisejsonala",
		"ie":         "utf-8",
		"fromsf":     "1",
		"word":       keyword,
		"pn":         0,
		"rn":         3,
		"gsm":        "3c",
		"searchtype": "0",
		"prefresh":   "undefined",
		"fromfilter": "0",
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, fmt.Sprint(v))
		req.URL.RawQuery = q.Encode()
	}
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("Referer", "https://m.baidu.com/sf/vsearch")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Mobile Safari/537.36")
	resp, err := searchHttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var data gson
	err = json.Unmarshal(bin, &data)
	if err != nil {
		log.Println(err, string(bin))
		return
	}
	linkData := data["linkData"].([]interface{})
	for _, item := range linkData {
		result = append(result, item.(gson)["objurl"].(string))
	}
	return
}

func checkFace(img []byte) (result []string) {
	link := upImg(img)
	formData := url.Values{}
	formData.Add("MsgId", fmt.Sprint(time.Now().Unix()*1000))
	formData.Add("CreateTime", fmt.Sprint(time.Now().Unix()))
	formData.Add("Content[imageUrl]", link)
	req, _ := http.NewRequest("POST", "https://kan.msxiaobing.com/Api/ImageAnalyze/Process?service=beauty&tid=", strings.NewReader(formData.Encode()))
	req.Header.Set("Referer", "https://kan.msxiaobing.com/ImageGame/Portal?task=beauty&feid=")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := faceHttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var data gson
	err = json.Unmarshal(bin, &data)
	if err != nil {
		log.Println(err, string(bin))
		return
	}
	log.Println(data, string(bin))
	return
}

func upImg(img []byte) (link string) {
	str := base64.StdEncoding.EncodeToString(img)
	req, _ := http.NewRequest("POST", "https://kan.msxiaobing.com/Api/Image/UploadBase64", strings.NewReader(str))
	req.Header.Set("Referer", "https://kan.msxiaobing.com/ImageGame/Portal?task=beauty&feid=")
	resp, err := faceHttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	var data gson
	err = json.Unmarshal(bin, &data)
	if err != nil {
		log.Println(err, string(bin))
		return
	}
	return data["Host"].(string) + data["Url"].(string)
}
