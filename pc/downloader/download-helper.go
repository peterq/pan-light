package downloader

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func redirectedLink(req *http.Request) (link string, err error) {
	c := http.Client{
		//Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Referer")
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "访问链接错误")
	}
	end := resp.Request.URL.String()
	resp.Body.Close()
	return end, nil
}

func downloadFileInfo(req *http.Request) (length int64, filename string, supportRange bool, err error) {

	c := http.Client{
		//Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Del("Referer")
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
	res, err := c.Do(req)
	if err != nil {
		err = errors.Wrap(err, "访问链接错误")
		return
	}
	if res.StatusCode != 200 {
		err = errors.Errorf("访问链接错误, http 状态码: %d", res.StatusCode)
		return
	}
	if cd, ok := res.Header["Content-Disposition"]; ok && len(cd) > 0 {
		if strings.IndexAny(cd[0], "attachment;filename=") != 0 {
			err = errors.New("不是文件链接")
			return
		}
		filename = strings.Trim(cd[0][len("attachment;filename="):], "\"")
	} else {
		err = errors.New("不是文件链接")
		return
	}
	if cl, ok := res.Header["Content-Length"]; ok && len(cl) > 0 {
		length, err = strconv.ParseInt(cl[0], 10, 64)
	} else {
		err = errors.New("不是文件链接")
		return
	}
	supportRange = strings.Contains(res.Header.Get("Accept-Ranges"), "bytes")
	return
}
