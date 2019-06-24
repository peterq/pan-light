package server_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/storage"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type gson = map[string]interface{}

var urlMap = map[string]string{
	"login-token":   "/api/pc/login-token",
	"login":         "/api/pc/login",
	"feedback":      "/api/pc/feedback",
	"refresh-token": "/api/pc/refresh-token",
	"share":         "/api/pc/share",
	"share-list":    "/api/pc/share/list",
	"share-hit":     "/api/pc/share/hit",
	"link-md5":      "/api/pc/link/md5",
}
var httpClient = http.Client{
	//Timeout: 15 * time.Second,
}

func makeRequest(name string, data map[string]interface{}) *http.Request {
	bin, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", dep.Env.ApiBase+urlMap[name], bytes.NewReader(bin))
	request.Header.Set("User-Agent", dep.Env.ClientUA)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "Bearer "+storage.UserState.Token)
	return request
}

func Call(name string, param map[string]interface{}) (result interface{}, err error) {
	res, err := httpClient.Do(makeRequest(name, param))
	if err != nil {
		return
	}
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var ret gson
	err = json.Unmarshal(bin, &ret)
	if err == nil {
		if !ret["success"].(bool) {
			return ret, errors.New(fmt.Sprint("api error(", ret["code"], "): ", ret["message"]))
		}
		result = ret["result"]
	} else {
		err = errors.Wrap(err, "json resp invalid")
		log.Println(string(bin))
	}
	return
}
