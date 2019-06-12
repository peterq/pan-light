package nickname

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/peterq/pan-light/server/cmd/cv"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type taskItem struct {
	book            string
	nickname        string
	images          []string
	faceCheckResult []struct {
		img   image.Image
		rects []image.Rectangle
	}
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
	go func() {
		for book, nicknames := range nicknameMap {
			for _, nickname := range nicknames {
				if fileExist(fmt.Sprintf("./data/avatar/result/%s/%s.jpg", book, nickname)) {
					continue
				}
				nicknameChan <- &taskItem{
					book:     book,
					nickname: nickname,
					images:   []string{},
					faceCheckResult: []struct {
						img   image.Image
						rects []image.Rectangle
					}{},
				}
			}
		}
		close(nicknameChan)
	}()
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
	savePath := fmt.Sprintf("./data/avatar/result/%s/%s.jpg", item.book, item.nickname)
	os.MkdirAll(path.Dir(savePath), os.ModePerm)
	for _, result := range item.faceCheckResult {
		if len(result.rects) != 1 {
			continue
		}
		// rect扩大2倍
		rect := result.rects[0]
		clip := rect
		clip.Min.X -= rect.Dx() / 2
		clip.Min.Y -= rect.Dy() / 2
		clip.Max.X += rect.Dx() / 2
		clip.Max.Y += rect.Dy() / 2
		if clip.Min.X < 0 {
			clip.Min.X = 0
		}
		if clip.Min.Y < 0 {
			clip.Min.Y = 0
		}
		if clip.Max.X > result.img.Bounds().Dx() {
			clip.Max.X = result.img.Bounds().Dx()
		}
		if clip.Max.Y > result.img.Bounds().Dy() {
			clip.Max.Y = result.img.Bounds().Dy()
		}
		type subImager interface {
			SubImage(image.Rectangle) image.Image
		}
		subImg := result.img.(subImager).SubImage(clip)
		buf := bytes.NewBuffer([]byte{})
		err := jpeg.Encode(buf, subImg, nil)
		if err != nil {
			continue
		}
		err = ioutil.WriteFile(savePath, buf.Bytes(), os.ModePerm)
		if err != nil {
			continue
		}
		log.Println("result", item.nickname, "success")
		return
	}
	log.Println("result", item.nickname, "fail")
	textAvatar(item.nickname, savePath)
}

func searchImageLoop(nicknameChan chan *taskItem, imageChan chan *taskItem, wg *sync.WaitGroup) {
	go func() {
		for item := range nicknameChan {
			log.Println("search item", item.book, item.nickname)
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
			i := 1
			for _, link := range item.images {
				log.Println("download", link)
				srcPath := fmt.Sprintf("./data/avatar/original/%s/%s_%d.jpg", item.book, item.nickname, i)
				markedPath := fmt.Sprintf("./data/avatar/marked/%s/%s_%d.jpg", item.book, item.nickname, i)
				img, err := downloadImg(srcPath, link)
				if err != nil {
					log.Println(err)
					continue
				}
				rect := cv.CheckFace(srcPath, markedPath)
				item.faceCheckResult = append(item.faceCheckResult, struct {
					img   image.Image
					rects []image.Rectangle
				}{img: img, rects: rect})
				i++
			}
			resultChan <- item
		}
		wg.Done()
	}()
	return
}

var searchHttpClient http.Client
var faceHttpClient http.Client
var font *truetype.Font

func avatarInit() {
	parseNicknameDoc()
	bin, err := ioutil.ReadFile("./data/font/cn.ttf")
	if err != nil {
		panic(err)
	}
	font, err = freetype.ParseFont(bin)
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

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func downloadImg(imgPath string, link string) (img image.Image, err error) {
	resp, err := http.Get(link)
	if err != nil {
		return
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	img, f, err := image.Decode(bytes.NewReader(bin))
	if err != nil {
		return
	}
	if f != "jpeg" {
		buf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(buf, img, nil)
		if err != nil {
			return
		}
		bin = buf.Bytes()
	}
	os.MkdirAll(path.Dir(imgPath), os.ModePerm)
	err = ioutil.WriteFile(imgPath, bin, os.ModePerm)
	return
}

func textAvatar(nickname, savePath string) {
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	rand.Seed(time.Now().UnixNano())
	r, g, b, _ := hls(rand.Intn(360))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}), image.ZP, draw.Src)

	str := []rune(nickname)
	if len(str) > 4 {
		str = []rune{str[0]}
	}
	var x, y, fontSize int
	ln1, ln2 := string(str), ""
	if len(str) == 1 {
		x, y, fontSize = 45, 700, img.Bounds().Dx()/3*2
	} else if len(str) == 2 {
		x, y, fontSize = 20, 700, img.Bounds().Dx()/2
	} else if len(str) == 3 {
		x, y, fontSize = 20, 700, img.Bounds().Dx()/3
	} else {
		ln1 = string(str[0]) + string(str[1])
		ln2 = string(str[2]) + string(str[3])
		x, y, fontSize = 60, 300, img.Bounds().Dx()/3
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)

	c.SetFontSize(float64(fontSize))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(color.White))
	// Draw the text.
	pt := freetype.Pt(x, x+int(c.PointToFixed(float64(y)))>>8)
	pt1, err := c.DrawString(ln1, pt)
	if err != nil {
		log.Println(err)
	}
	pt.Y += (pt1.X-pt.X)/2 + c.PointToFixed(20)
	_, err = c.DrawString(ln2, pt)
	if err != nil {
		log.Println(err)
	}

	os.MkdirAll(path.Dir(savePath), os.ModePerm)
	buf := bytes.NewBuffer([]byte{})
	jpeg.Encode(buf, img, nil)
	ioutil.WriteFile(savePath, buf.Bytes(), os.ModePerm)
}

func hls(H int) (r, g, b, a int) {
	S, V := 255, 255
	// Direct implementation of the graph in this image:
	// https://en.wikipedia.org/wiki/HSL_and_HSV#/media/File:HSV-RGB-comparison.svg
	max := V
	min := V * (255 - S)

	H %= 360
	segment := H / 60
	offset := H % 60
	mid := ((max - min) * offset) / 60

	//log.Println(H, max, min, mid)
	switch segment {
	case 0:
		return max, min + mid, min, 0xff
	case 1:
		return max - mid, max, min, 0xff
	case 2:
		return min, max, min + mid, 0xff
	case 3:
		return min, max - mid, max, 0xff
	case 4:
		return min + mid, min, max, 0xff
	case 5:
		return max, min, max - mid, 0xff
	}

	return 0, 0, 0, 0xff
}
