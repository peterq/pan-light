package storage

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/peterq/pan-light/pc/dep"
	"github.com/peterq/pan-light/pc/util"
	"io/ioutil"
	"log"
	"time"
)

var dataFilePath = dep.DataPath("pan-light.data")

var UserState *State
var Global = &GlobalData{}

func init() {
	dep.OnInit(doInit)
	dep.OnClose(func() {
		log.Println("退出, 数据存盘...")
	})
}

func firstStart() {
	Global = &GlobalData{
		UserStateMap: map[string]*State{
			"default": firstLogin("default"),
		},
		CurrentUser: "default",
	}
}

func firstLogin(username string) *State {
	return &State{
		Token: "",
		Settings: &StateSetting{
			DownloadSegSize:   2,
			DownloadCoroutine: 128,
		},
		PanCookie:   []*Cookies{},
		UserStorage: map[string]string{},
		Username:    username,
		Uk:          "",
		Logout:      false,
	}
}

func doInit() {
	defer writeToDiskTask()
	if util.First(util.PathExists(dataFilePath)).(bool) {
		err := proto.Unmarshal(util.First(ioutil.ReadFile(dataFilePath)).([]byte), Global)
		if err != nil {
			panic(err)
		}
	} else {
		firstStart()
	}
	Global.UserStateMap["default"] = firstLogin("default")
	UserState = Global.UserStateMap[Global.CurrentUser]
}

func OnLogin(username string) {
	if username == Global.CurrentUser {
		return
	}
	Global.CurrentUser = username
	state, ok := Global.UserStateMap[username]
	if ok {
		UserState = state
	} else {
		UserState = firstLogin(username)
		Global.UserStateMap[username] = UserState
	}
	UserState.Logout = false
}

func writeToDiskTask() {
	var lastSerialize []byte
	fn := func() {
		if Global == nil {
			return
		}
		bin, err := proto.Marshal(Global)
		if err != nil {
			dep.Fatal(err.Error())
		}
		if !bytes.Equal(lastSerialize, bin) {
			err = ioutil.WriteFile(dataFilePath, bin, 0655)
			if err != nil {
				dep.Fatal(err.Error())
			}
			lastSerialize = bin
		}
	}
	dep.OnClose(fn)
	go func() {
		for range time.Tick(10 * time.Second) {
			fn()
		}
	}()
}

func UserStorageSet(k, v string) {
	UserState.UserStorage[k] = v
}

func UserStorageGet(k string) string {
	v, ok := UserState.UserStorage[k]
	if ok {
		return v
	}
	return ""
}
